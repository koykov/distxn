package distxn

type Jobs struct {
	buf []Job
}

func (j *Jobs) AddJob(job Job) {
	j.buf = append(j.buf, job)
}

func (j *Jobs) Jobs() []Job {
	return j.buf
}

func (j *Jobs) Reset() {
	j.buf = j.buf[:0]
}
