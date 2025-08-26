# LissanAI Backend

Complete authentication system with email notifications for password reset functionality.

## 📧 Email Configuration

The system now supports sending password reset emails via SMTP. Configure the following environment variables:

### Required Email Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `SMTP_HOST` | SMTP server hostname | `smtp.gmail.com` |
| `SMTP_PORT` | SMTP server port | `587` |
| `SMTP_USERNAME` | SMTP username/email | `your-email@gmail.com` |
| `SMTP_PASSWORD` | SMTP password/app password | `your-app-password` |
| `FROM_EMAIL` | From email address | `noreply@lissanai.com` |
| `FRONTEND_URL` | Frontend application URL | `http://localhost:3000` |

### Email Setup Options

#### Option 1: Gmail SMTP (Recommended for Development)
```env
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-gmail@gmail.com
SMTP_PASSWORD=your-app-password
FROM_EMAIL=your-gmail@gmail.com
FRONTEND_URL=http://localhost:3000
```

**Gmail Setup Steps:**
1. Enable 2-factor authentication on your Gmail account
2. Generate an App Password: Google Account → Security → App passwords
3. Use the generated app password (not your regular password)

#### Option 2: Development Mode (No Email Sending)
If SMTP credentials are not configured, the system will:
- Print reset links to console instead of sending emails
- Continue working normally for development/testing

#### Option 3: Production SMTP Services
- **SendGrid**: `smtp.sendgrid.net:587`
- **Mailgun**: `smtp.mailgun.org:587`
- **AWS SES**: `email-smtp.region.amazonaws.com:587`

### Password Reset Email Features

- **HTML Email Template**: Professional, responsive design
- **Security Warnings**: Clear instructions about link expiry and security
- **1-Hour Expiry**: Reset tokens expire automatically
- **Branded Design**: LissanAI themed email template
- **Fallback Text**: Plain text version for accessibility

### Testing Email Functionality

```bash
# Test forgot password endpoint
curl -X POST http://localhost:8080/api/v1/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com"}'

# Check console output if SMTP not configured
# Or check your email inbox if SMTP is configured
```

## 🔧 Environment Variables

| Variable | Description | Default/Example |
|----------|-------------|-----------------|
| `MONGODB_URI` | MongoDB connection string | `mongodb://localhost:27017/lissanai` |
| `JWT_SECRET` | Secret key for JWT tokens | `your-secret-key` |
| `SMTP_HOST` | SMTP server hostname | `smtp.gmail.com` |
| `SMTP_PORT` | SMTP server port | `587` |
| `SMTP_USERNAME` | SMTP username/email | `your-email@gmail.com` |
| `SMTP_PASSWORD` | SMTP password/app password | `your-app-password` |
| `FROM_EMAIL` | From email address | `noreply@lissanai.com` |
| `FRONTEND_URL` | Frontend application URL | `http://localhost:3000` |