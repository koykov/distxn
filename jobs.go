package distnx

type jobs struct {
    buf []Job
}

func (j *jobs) AddJob(job Job) {
    j.buf = append(j.buf, job)
}
