package gitlab

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	goGitLab "github.com/xanzy/go-gitlab"
)

const gitlabTokenEnv = "GITLAB_TOKEN"
const jobIDFile = "jobs.txt"
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

func getPipelineIDs(page int) []int {

	cmd := exec.Command("sh", "pipelines.sh", fmt.Sprintf("%d", page))
	fmt.Println("INDEX", page)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(jobIDFile)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var jobIDs []int
	for scanner.Scan() {
		id, err := strconv.Atoi(scanner.Text())
		if err != nil {
			continue
		}
		jobIDs = append(jobIDs, id)
	}

	return jobIDs
}

// PagesFlush cleans up the pipelines and their jobs in the CI.
func PagesFlush() string {
	log.Println("Flushing CI for GitLab Pages...")

	projID, git := PagesClient()
	flushed := 0
	pipelinePageCount := 1
	pipelines := []int{0}
	for len(pipelines) > 0 {
		pipelines = getPipelineIDs(pipelinePageCount)
		fmt.Println(pipelines)
		for _, pipeline := range pipelines {
			jobs, _, err := git.Jobs.ListPipelineJobs(
				projID,
				pipeline,
				&goGitLab.ListJobsOptions{})
			if err != nil {
				continue
			}

			for _, job := range jobs {
				fmt.Println(job.ID)
				if _, _, err := git.Jobs.EraseJob(projID, job.ID); err != nil {
					// log.Println(err)
				} else {
					flushed++
				}
			}
			log.Printf("Flushed job traces for pipeline %d.\n", pipeline)
		}

		pipelinePageCount++
	}

	return fmt.Sprintf("Flushed %d jobs.", flushed)
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
