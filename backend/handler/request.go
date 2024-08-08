package handler

type getUserRequest struct {
	ID string `validate:"required,id"`
}

func (r *getUserRequest) bind(v *Validator) error {
	if err := v.Validate(r); err != nil {
		return err
	}

	return nil
}

type insertUserRequest struct {
	FirstName string `validate:"required,alpha,min=2,max=48"`
	LastName  string `validate:"required,alpha,min=2,max=48"`
	Email     string `validate:"required,email"`
	Password  string `validate:"required,min=7,max=256"`
}

func (r *insertUserRequest) bind(v *Validator) error {
	if err := v.Validate(r); err != nil {
		return err
	}

	return nil
}

type updateUserRequest struct {
	ID        string `validate:"required,id"`
	FirstName string `validate:"omitempty,alpha,min=2,max=48"`
	LastName  string `validate:"omitempty,alpha,min=2,max=48"`
	Email     string `validate:"omitempty,email"`
	Password  string `validate:"omitempty,min=7,max=256"`
}

func (r *updateUserRequest) bind(v *Validator) error {
	if err := v.Validate(r); err != nil {
		return err
	}

	return nil
}
