package gitlab

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	goGitLab "github.com/xanzy/go-gitlab"
)

const gitlabTokenEnv = "GITLAB_TOKEN"
const projIDEnv = "PROJECT_ID"
const websitePipelinesBase = "https://gitlab.com/utmist/utmist.gitlab.io/pipelines"

// PagesClient returns GitLab client for the website.
func PagesClient() (string, *goGitLab.Client) {
	// Look for the two environment variables.
	token, exists := os.LookupEnv(gitlabTokenEnv)
	if !exists {
		reply := "No discord bot token found."
		log.Println(reply)
		return "", nil
	}
	projID, exists := os.LookupEnv(projIDEnv)
	if !exists {
		reply := "No GitLab project ID found."
		log.Println(reply)
		return "", nil
	}

	// Set up GL Client
	git := goGitLab.NewClient(nil, token)
	git.SetBaseURL("https://gitlab.com/api/v4")

	return projID, git
}

// PagesFlush cleans up the pipelines and their jobs in the CI.
func PagesFlush() string {
	log.Println("Flushing CI for GitLab Pages...")

	projID, git := PagesClient()
	listOpts := &goGitLab.ListProjectPipelinesOptions{
		Status: goGitLab.BuildState(goGitLab.Success),
	}
	pipelines, _, err := git.Pipelines.ListProjectPipelines(projID, listOpts)
	if err != nil || len(pipelines) == 0 {
		log.Fatal(err)
	}

	flushed := []string{}
	skipped := []string{}
	for _, pl := range pipelines {
		jobs, _, err := git.Jobs.ListPipelineJobs(projID, pl.ID,
			&goGitLab.ListJobsOptions{})
		if err != nil {
			continue
		}

		for _, job := range jobs {
			if job.FinishedAt.Before(time.Now().AddDate(0, 0, -1)) {
				git.Jobs.EraseJob(projID, job.ID)
				flushed = append(flushed, fmt.Sprintf("%d", job.ID))
				continue
			}
			log.Printf("Skipping %d; it's within last 24 h.", job.ID)
			skipped = append(skipped, fmt.Sprintf("%d", job.ID))
		}
	}

	return fmt.Sprintf("Skipped {%s}\nFlushed {%s}.",
		strings.Join(skipped, ", "),
		strings.Join(flushed, ", "))
}

// PagesUpdate reruns last successful pipeline on master for utmist.gitlab.io.
func PagesUpdate() string {
	log.Println("Running CI for GitLab Pages...")

	projID, git := PagesClient()

	// Search for the most recent successful pipeline on master.
	listOpts := &goGitLab.ListProjectPipelinesOptions{
		Status: goGitLab.BuildState(goGitLab.Success),
		Ref:    goGitLab.String("master"),
	}
	pipelines, _, err := git.Pipelines.ListProjectPipelines(projID, listOpts)
	if err != nil || len(pipelines) == 0 {
		log.Fatal(err)
	}
	pipeline := pipelines[0]
	log.Printf("Found successful pipeline: %d\n", pipeline.ID)

	// Get variables of last successful pipeline on master.
	vars, _, err := git.Pipelines.GetPipelineVariables(projID, pipeline.ID)
	if err != nil {
		log.Fatal(err)
	}

	// Use variables to run new pipeline.
	opt := &goGitLab.CreatePipelineOptions{
		Ref:       goGitLab.String("master"),
		Variables: vars,
	}
	newPipeline, _, err := git.Pipelines.CreatePipeline(projID, opt)
	if err != nil {
		log.Fatal(err)
	}

	// Create and return successful reply.
	reply := fmt.Sprintf("Successfully rerun pipeline: %d as %d\nSee %s.\n",
		pipeline.ID,
		newPipeline.ID,
		fmt.Sprintf("%s/%d", websitePipelinesBase, newPipeline.ID))

	log.Print(reply)
	return reply
}
