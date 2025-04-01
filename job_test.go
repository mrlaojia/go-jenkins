package jenkinsdk

import "testing"

/**
* @Author: jack.walker
* @File: job_test.go
* @CreateDate: 2025/4/1 14:23
* @Version: 1.0.0
* @Description:
 */

const (
	jobConfig = `<?xml version='1.1' encoding='UTF-8'?>
        <project>
          <actions/>
          <description>i am new-s job description</description>
          <keepDependencies>false</keepDependencies>
          <properties/>
          <scm class="hudson.scm.NullSCM"/>
          <canRoam>true</canRoam>
          <disabled>false</disabled>
          <blockBuildWhenDownstreamBuilding>false</blockBuildWhenDownstreamBuilding>
          <blockBuildWhenUpstreamBuilding>false</blockBuildWhenUpstreamBuilding>
          <triggers/>
          <concurrentBuild>false</concurrentBuild>
          <builders>
            <hudson.tasks.Shell>
              <command>echo &quot;dev job--5&quot;</command>
              <configuredLocalRules/>
            </hudson.tasks.Shell>
          </builders>
          <publishers/>
          <buildWrappers/>
        </project>`
)

func TestJenkinsSdk_CreateJob(t *testing.T) {
	j := NewJenkinsSdk("http://172.19.89.76:48080/", "wkj", "11900ac516d0ac841dfcdab7ed042b9fcb")

	job := &JenkinsJob{
		Name:      "job-create_2",
		ConfigXml: jobConfig,
	}

	err := j.CreateJob(job)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("create job success")
	}
}

func TestJenkinsSdk_CopyJob(t *testing.T) {
	j := NewJenkinsSdk("http://172.19.89.76:48080/", "wkj", "11900ac516d0ac841dfcdab7ed042b9fcb")

	fromJob := &JenkinsJob{
		Name: "job-create_2",
		//Parent: []string{"dev1"},
	}

	job := &JenkinsJob{
		Name: "job-copy_5",
	}

	err := j.CopyJob(job, fromJob)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("copy job success")
	}

	err = j.DisableJob(job)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("disable job success")
	}

	err = j.EnableJob(job)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("enable job success")
	}
}

func TestJenkinsSdk_GetJob(t *testing.T) {
	j := NewJenkinsSdk("http://172.19.89.76:48080/", "wkj", "11900ac516d0ac841dfcdab7ed042b9fcb")

	job := &JenkinsJob{
		Name: "job-copy_1",
	}

	config, err := j.GetJob(job)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("get job success:\n%v", string(config))
	}
}

func TestJenkinsSdk_DeleteJob(t *testing.T) {
	j := NewJenkinsSdk("http://172.19.89.76:48080/", "wkj", "11900ac516d0ac841dfcdab7ed042b9fcb")

	job := &JenkinsJob{
		Name:   "dev_job-1",
		Parent: []string{"dev1"},
	}

	err := j.DeleteJob(job)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("delete job success")
	}
}

func TestJenkinsSdk_UpdateJobDescription(t *testing.T) {
	j := NewJenkinsSdk("http://172.19.89.76:48080/", "wkj", "11900ac516d0ac841dfcdab7ed042b9fcb")

	job := &JenkinsJob{
		Name: "job-copy_2",
	}

	err := j.UpdateJobDescription(job, "i am new-s job description")
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("update job description success")
	}

}

func TestJenkinsSdk_UpdateJob(t *testing.T) {
	j := NewJenkinsSdk("http://172.19.89.76:48080/", "wkj", "11900ac516d0ac841dfcdab7ed042b9fcb")

	job := &JenkinsJob{
		Name:      "job-copy_5",
		ConfigXml: jobConfig,
	}

	err := j.UpdateJob(job)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("update job success")
	}
}
