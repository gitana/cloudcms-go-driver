package cloudcms

import (
	"fmt"
	"time"
)

func (session *CloudCmsSession) ReadPlatform() (JsonObject, error) {
	return session.Get("/", nil)
}

// def wait_for_job_completion(self, jobId):

// # Use with caution
// while True:
// 	job = self.read_job(jobId)

// 	if job.data['state'] == 'FINISHED':
// 		return job
// 	elif job.data['state'] == 'ERROR':
// 		raise JobError(jobId)
// 	else:
// 		time.sleep(1)

func (session *CloudCmsSession) ReadJob(jobId string) (JsonObject, error) {
	return session.Get(fmt.Sprintf("/jobs/%s", jobId), nil)
}

func (session *CloudCmsSession) WaitForJob(jobId string) error {
	for {
		job, err := session.ReadJob(jobId)
		if err != nil {
			return err
		}
		if job.GetString("state") == "FINISHED" {
			return nil
		}
		if job.GetString("state") == "ERROR" {
			return fmt.Errorf("job failed: %s", jobId)
		}

		time.Sleep(1 * time.Second)
	}
}
