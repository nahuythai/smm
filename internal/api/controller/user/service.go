package user

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/url"
	"smm/internal/database/models"
	"smm/internal/database/queries"
	"smm/pkg/constants"
	"smm/pkg/jwt"
	"smm/pkg/response"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct {
}

type serviceInterface interface {
	sessionValidate(ctx context.Context, sessionToken string, sessionType int) (*models.Session, error)
	createVerifyEmailTemplate(username, token string) (string, error)
	verifyEmailSuccessTemplate() (string, error)
	verifyEmailFailTemplate() (string, error)
}

func NewService() serviceInterface {
	return new(service)
}

func (s *service) createVerifyEmailTemplate(username, token string) (string, error) {
	emailVerificationTemplate := `
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>Email Verification</title>
		<style>
			body {
				font-family: Arial, sans-serif;
				background-color: #f2f2f2;
				margin: 0;
				padding: 0;
			}
	
			.container {
				max-width: 600px;
				margin: 0 auto;
				padding: 20px;
			}
	
			.logo {
				text-align: center;
				margin-bottom: 20px;
			}
	
			.logo img {
				max-width: 150px;
			}
	
			.content {
				background-color: #ffffff;
				padding: 30px;
				border-radius: 5px;
				box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
			}
	
			h2 {
				color: #333333;
				font-size: 24px;
				margin-top: 0;
			}
	
			p {
				color: #555555;
				font-size: 16px;
				line-height: 1.5;
				margin: 0 0 20px;
			}
	
			.btn {
				display: inline-block;
				background-color: #007bff;
				color: #ffffff;
				text-decoration: none;
				padding: 10px 20px;
				border-radius: 3px;
			}
	
			.btn:hover {
				background-color: #0056b3;
			}
	
			.footer {
				margin-top: 40px;
				text-align: center;
				color: #999999;
			}
	
			.footer p {
				margin: 0;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<div class="logo">
				<img src="https://daily.like1000d.com/assets/uploads/logo/logo.png" alt="Logo">
			</div>
			<div class="content">
				<h2>Email Verification</h2>
				<p>Dear {{.Name}},</p>
				<p>Thank you for signing up. Please click the button below to verify your email address:</p>
				<p>
					<a class="btn" href="{{.VerificationLink}}">Verify Email</a>
				</p>
				<p>If you did not sign up for this account, you can safely ignore this email.</p>
			</div>
			<div class="footer">
				<p>Best regards, <br>1000like</p>
			</div>
		</div>
	</body>
	</html>
	`

	// Create the verification link
	verificationLink := fmt.Sprintf("http://%s/api/v1/users/verify-email?token=%s", cfg.ServerDomain, url.QueryEscape(token))

	// Create a template object
	tmpl := template.Must(template.New("emailVerification").Parse(emailVerificationTemplate))

	// Prepare the data for rendering the template
	data := struct {
		Name             string
		VerificationLink string
	}{
		Name:             username, // Replace with the recipient's name
		VerificationLink: verificationLink,
	}

	// Render the template to stdout
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (s *service) verifyEmailFailTemplate() (string, error) {
	emailVerificationSuccessTemplate := `
	<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Email Verification Successful</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f2f2f2;
            margin: 0;
            padding: 0;
        }

        .container {
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }

        .logo {
            text-align: center;
            margin-bottom: 20px;
        }

        .logo img {
            max-width: 150px;
        }

        .content {
            background-color: #ffffff;
            padding: 30px;
            border-radius: 5px;
            box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
        }

        h2 {
            color: #333333;
            font-size: 24px;
            margin-top: 0;
        }

        p {
            color: #555555;
            font-size: 16px;
            line-height: 1.5;
            margin: 0 0 20px;
        }

        .success-msg {
            color: #28a745;
            font-weight: bold;
        }

        .footer {
            margin-top: 40px;
            text-align: center;
            color: #999999;
        }

        .footer p {
            margin: 0;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="logo">
            <img src="https://daily.like1000d.com/assets/uploads/logo/logo.png" alt="Logo">
        </div>
        <div class="content">
		<h2>Email Verification Failed</h2>
		<p>We're sorry, but your email verification has failed. Please try again or contact our support team for assistance.</p>
        </div>
        <div class="footer">
            <p>Best regards, <br>1000Like</p>
        </div>
    </div>
</body>
</html>
	`
	return emailVerificationSuccessTemplate, nil
}

func (s *service) verifyEmailSuccessTemplate() (string, error) {
	emailVerificationSuccessTemplate := `
	<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Email Verification Successful</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f2f2f2;
            margin: 0;
            padding: 0;
        }

        .container {
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }

        .logo {
            text-align: center;
            margin-bottom: 20px;
        }

        .logo img {
            max-width: 150px;
        }

        .content {
            background-color: #ffffff;
            padding: 30px;
            border-radius: 5px;
            box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
        }

        h2 {
            color: #333333;
            font-size: 24px;
            margin-top: 0;
        }

        p {
            color: #555555;
            font-size: 16px;
            line-height: 1.5;
            margin: 0 0 20px;
        }

        .success-msg {
            color: #28a745;
            font-weight: bold;
        }

        .footer {
            margin-top: 40px;
            text-align: center;
            color: #999999;
        }

        .footer p {
            margin: 0;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="logo">
            <img src="https://daily.like1000d.com/assets/uploads/logo/logo.png" alt="Logo">
        </div>
        <div class="content">
            <h2>Email Verification Successful</h2>
            <p>Your email address has been successfully verified. You can now access all the features of our platform.</p>
            <p class="success-msg">Thank you for verifying your email!</p>
        </div>
        <div class="footer">
            <p>Best regards, <br>1000Like</p>
        </div>
    </div>
</body>
</html>
	`
	return emailVerificationSuccessTemplate, nil
}

func (s *service) sessionValidate(ctx context.Context, sessionToken string, sessionType int) (*models.Session, error) {
	payload, err := jwt.GetGlobal().ValidateToken(sessionToken)
	if err != nil {
		logger.Error().Err(err).Caller().Str("func", "sessionValidate").Str("funcInline", "jwt.GetGlobal().ValidateToken").Msg("user-controller")
		return nil, response.NewError(fiber.StatusUnauthorized, response.Option{Code: constants.ErrCodeTokenWrong, Data: "missing or wrong token"})
	}
	if payload.Type != constants.TokenTypeSession {
		return nil, response.NewError(fiber.StatusUnauthorized, response.Option{Code: constants.ErrCodeTokenWrong, Data: "missing or wrong token"})
	}
	id, _ := primitive.ObjectIDFromHex(payload.ID)
	session, err := queries.NewSession(ctx).GetById(id)
	if err != nil {
		if e, ok := err.(*response.Option); ok {
			if e.Code == constants.ErrCodeSessionNotFound {
				return nil, response.NewError(fiber.StatusUnauthorized, response.Option{Code: constants.ErrCodeSessionNotFound})
			}
		}
	}
	if session.Type != sessionType {
		return nil, response.NewError(fiber.StatusUnauthorized, response.Option{Code: constants.ErrCodeTokenWrong, Data: "missing or wrong token"})
	}
	return session, nil
}
