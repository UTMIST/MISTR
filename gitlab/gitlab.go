package gitlab

import (
	"fmt"
	"log"
	"os"

	goGitLab "github.com/xanzy/go-gitlab"
)

const gitlabTokenEnv = "GITLAB_TOKEN"
const projectIDEnv = "PROJECT_ID"
const websitePipelinesBase = "https://gitlab.com/utmist/utmist.gitlab.io/pipelines"

// PagesUpdate reruns last successful pipeline on master for utmist.gitlab.io.
func PagesUpdate() string {
	log.Println("Running CI for GitLab Pages...")

	// Look for the two environment variables.
	token, exists := os.LookupEnv(gitlabTokenEnv)
	if !exists {
		reply := "No discord bot token found."
		log.Println(reply)
		return reply
	}
	projectID, exists := os.LookupEnv(projectIDEnv)
	if !exists {
		reply := "No GitLab project ID found."
		log.Println(reply)
		return reply
	}

	// Set up GL Client
	git := goGitLab.NewClient(nil, token)
	git.SetBaseURL("https://gitlab.com/api/v4")

	// Search for the most recent successful pipeline on master.
	listOpts := &goGitLab.ListProjectPipelinesOptions{
		Status: goGitLab.BuildState(goGitLab.Success),
		Ref:    goGitLab.String("master"),
	}
	pipelines, _, err := git.Pipelines.ListProjectPipelines(projectID, listOpts)
	if err != nil || len(pipelines) == 0 {
		log.Fatal(err)
	}
	pipeline := pipelines[0]
	log.Printf("Found successful pipeline: %d\n", pipeline.ID)

	// Get variables of last successful pipeline on master.
	vars, _, err := git.Pipelines.GetPipelineVariables(projectID, pipeline.ID)
	if err != nil {
		log.Fatal(err)
	}

	// Use variables to run new pipeline.
	opt := &goGitLab.CreatePipelineOptions{
		Ref:       goGitLab.String("master"),
		Variables: vars,
	}
	newPipeline, _, err := git.Pipelines.CreatePipeline(projectID, opt)
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
