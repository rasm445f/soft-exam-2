package domain

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"regexp"

	"github.com/rasm445f/soft-exam-2/db/generated"
	"github.com/rasm445f/soft-exam-2/mailer"
)

type CustomerDomain struct {
	Queries *generated.Queries
}

// Define individual regex patterns
var (
	minLengthRegex   = regexp.MustCompile(`^.{8,}$`)
	lowercaseRegex   = regexp.MustCompile(`[a-z]`)
	uppercaseRegex   = regexp.MustCompile(`[A-Z]`)
	numberRegex      = regexp.MustCompile(`\d`)
	specialCharRegex = regexp.MustCompile(`[@$!%*?&]`)
)

// ValidatePassword checks if the password meets all complexity requirements
func ValidatePassword(password string) error {
	if !minLengthRegex.MatchString(password) {
		return errors.New("password must be at least 8 characters long")
	}
	if !lowercaseRegex.MatchString(password) {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !uppercaseRegex.MatchString(password) {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !numberRegex.MatchString(password) {
		return errors.New("password must contain at least one number")
	}
	if !specialCharRegex.MatchString(password) {
		return errors.New("password must contain at least one special character (@$!%*?&)")
	}
	return nil
}

// ValidateEmail checks if the email is in a valid format.
func ValidateEmail(email string) error {
	// Define a regex pattern for a valid email address
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Compile the regex
	re := regexp.MustCompile(emailRegex)

	// Match the email against the regex
	if !re.MatchString(email) {
		return errors.New("invalid email format")
	}

	return nil
}

func NewCustomerDomain(queries *generated.Queries) *CustomerDomain {
	return &CustomerDomain{Queries: queries}
}

func (cd *CustomerDomain) GetAllCustomersDomain(ctx context.Context) ([]generated.Customer, error) {
	return cd.Queries.GetAllCustomers(ctx)
}

func (cd *CustomerDomain) GetCustomerByIdDomain(ctx context.Context, id int32) (generated.Customer, error) {
	return cd.Queries.GetCustomerByID(ctx, id)
}

func (cd *CustomerDomain) DeleteCustomerDomain(ctx context.Context, id int32) error {
	return cd.Queries.DeleteCustomer(ctx, id)
}

func (cs *CustomerDomain) CreateCustomerDomain(ctx context.Context, customerParams generated.CreateCustomerParams) (*generated.Customer, error) {
	if *customerParams.Name == "" || *customerParams.Email == "" || *customerParams.Password == "" {
		return nil, errors.New("all required fields must be filled")
	}

	if err := ValidatePassword(*customerParams.Password); err != nil {
		return nil, err
	}

	createdCustomer, err := cs.Queries.CreateCustomer(ctx, customerParams)
	if err != nil {
		return nil, err
	}

	subject := "Welcome to MTOGO, " + *customerParams.Name + "!"

	body := `
    <html>
        <body style="font-family: Arial, sans-serif; color: #333;">
            <h1 style="color: #4CAF50;">Welcome to MTOGO!</h1>
            <p>Hi ` + *customerParams.Name + `,</p>
            <p>We're thrilled to have you join our community! Thank you for signing up with MTOGO.</p>
            
            <p>Here’s what you can look forward to as a new member:</p>
            <ul>
                <li><strong>Personalized Experience:</strong> Tailored recommendations and insights just for you.</li>
                <li><strong>Exclusive Access:</strong> Enjoy early access to new features and updates.</li>
                <li><strong>Dedicated Support:</strong> Our team is here to assist you whenever you need.</li>
            </ul>

            <p>To get started, simply log in and explore. We’re here to make sure you have a seamless experience, so don’t hesitate to reach out if you have any questions.</p>

            <p style="margin-top: 30px;">Cheers,</p>
            <p>The [Your Service Name] Team</p>
            <footer style="margin-top: 20px; font-size: 0.9em; color: #666;">
                <hr>
                <p>If you did not sign up for this account, please ignore this email.</p>
            </footer>
        </body>
    </html>
`

	err = mailer.SendMailWithGomail(*customerParams.Email, subject, body)
	if err != nil {
		log.Println("Failed to send email:", err)
	}

	customer := &generated.Customer{
		ID:          createdCustomer.ID,
		Name:        createdCustomer.Name,
		Email:       createdCustomer.Email,
		Password:    createdCustomer.Password,
		Phonenumber: createdCustomer.Phonenumber,
		Address:     createdCustomer.Address,
	}

	return customer, nil
}

func (cd *CustomerDomain) UpdateCustomerDomain(ctx context.Context, customerParams generated.UpdateCustomerParams) error {
	// Validate optional fields if provided
	if customerParams.Password != nil {
		if err := ValidatePassword(*customerParams.Password); err != nil {
			return err
		}
	}

	if customerParams.Email != nil {
		if err := ValidateEmail(*customerParams.Email); err != nil {
			return err
		}
	}

	// Perform the update in the database
	err := cd.Queries.UpdateCustomer(ctx, customerParams)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("customer not found")
		}
		return err
	}

	return nil
}
