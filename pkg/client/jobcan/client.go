type JobcanClient interface {
    Login(email, password string) error
    Adit() (*AditResult, error)
	GetAditErrors() (AditErrors, error)
}

