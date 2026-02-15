package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

// SecurityAuditEvent represents a security event to be logged
type SecurityAuditEvent struct {
	Timestamp time.Time              `json:"timestamp"`
	EventType string                 `json:"event_type"`
	UserID    uint                   `json:"user_id,omitempty"`
	Email     string                 `json:"email,omitempty"`
	IPAddress string                 `json:"ip_address,omitempty"`
	Action    string                 `json:"action"`
	Status    string                 `json:"status"` // "success" or "failure"
	Reason    string                 `json:"reason,omitempty"`
	UserAgent string                 `json:"user_agent,omitempty"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

// SecurityAuditor logs security-related events
type SecurityAuditor struct {
	logger *log.Logger
}

// NewSecurityAuditor creates a new security auditor
func NewSecurityAuditor() *SecurityAuditor {
	// Log to both stdout and a file
	logFile := os.Getenv("SECURITY_LOG_FILE")
	if logFile == "" {
		logFile = "/tmp/security-audit.log"
	}

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		fmt.Printf("Warning: Could not open security log file: %v\n", err)
		return &SecurityAuditor{
			logger: log.New(os.Stdout, "[SECURITY] ", log.LstdFlags),
		}
	}

	return &SecurityAuditor{
		logger: log.New(file, "[SECURITY] ", log.LstdFlags),
	}
}

// LogAuthenticationAttempt logs a login attempt (success or failure)
func (a *SecurityAuditor) LogAuthenticationAttempt(email, ipAddress, userAgent string, success bool, reason string, userID uint) {
	status := "success"
	eventType := "SUCCESSFUL_LOGIN"
	if !success {
		status = "failure"
		eventType = "FAILED_LOGIN"
	}

	event := SecurityAuditEvent{
		Timestamp: time.Now(),
		EventType: eventType,
		UserID:    userID,
		Email:     email,
		IPAddress: ipAddress,
		Action:    "login_attempt",
		Status:    status,
		Reason:    reason,
		UserAgent: userAgent,
	}

	a.logEvent(event)
}

// LogRegistration logs a user registration attempt
func (a *SecurityAuditor) LogRegistration(email, ipAddress, userAgent string, success bool, reason string, userID uint) {
	status := "success"
	eventType := "SUCCESSFUL_REGISTRATION"
	if !success {
		status = "failure"
		eventType = "FAILED_REGISTRATION"
	}

	event := SecurityAuditEvent{
		Timestamp: time.Now(),
		EventType: eventType,
		UserID:    userID,
		Email:     email,
		IPAddress: ipAddress,
		Action:    "registration",
		Status:    status,
		Reason:    reason,
		UserAgent: userAgent,
	}

	a.logEvent(event)
}

// LogLogout logs a user logout
func (a *SecurityAuditor) LogLogout(userID uint, email, ipAddress string) {
	event := SecurityAuditEvent{
		Timestamp: time.Now(),
		EventType: "LOGOUT",
		UserID:    userID,
		Email:     email,
		IPAddress: ipAddress,
		Action:    "logout",
		Status:    "success",
	}

	a.logEvent(event)
}

// LogUnauthorizedAccess logs unauthorized access attempts
func (a *SecurityAuditor) LogUnauthorizedAccess(requiredResource string, userID uint, ipAddress string, userAgent string) {
	event := SecurityAuditEvent{
		Timestamp: time.Now(),
		EventType: "UNAUTHORIZED_ACCESS",
		UserID:    userID,
		IPAddress: ipAddress,
		Action:    "access_attempt",
		Status:    "failure",
		Reason:    fmt.Sprintf("Unauthorized attempt to access: %s", requiredResource),
		UserAgent: userAgent,
	}

	a.logEvent(event)
}

// LogRateLimitExceeded logs when rate limit is exceeded
func (a *SecurityAuditor) LogRateLimitExceeded(endpoint string, ipAddress string, reason string) {
	event := SecurityAuditEvent{
		Timestamp: time.Now(),
		EventType: "RATE_LIMIT_EXCEEDED",
		IPAddress: ipAddress,
		Action:    "rate_limit_violation",
		Status:    "blocked",
		Reason:    fmt.Sprintf("Rate limit exceeded on %s: %s", endpoint, reason),
	}

	a.logEvent(event)
}

// LogSuspiciousActivity logs any suspicious activity
func (a *SecurityAuditor) LogSuspiciousActivity(activityType string, ipAddress string, userID uint, details map[string]interface{}) {
	event := SecurityAuditEvent{
		Timestamp: time.Now(),
		EventType: "SUSPICIOUS_ACTIVITY",
		UserID:    userID,
		IPAddress: ipAddress,
		Action:    activityType,
		Status:    "flagged",
		Details:   details,
	}

	a.logEvent(event)
}

// logEvent writes the event to the log
func (a *SecurityAuditor) logEvent(event SecurityAuditEvent) {
	jsonData, err := json.Marshal(event)
	if err != nil {
		a.logger.Printf("Error marshaling event: %v", err)
		return
	}

	a.logger.Printf("%s", string(jsonData))
}

// GetSecurityAuditor returns a singleton instance of SecurityAuditor
var globalAuditor *SecurityAuditor

// InitSecurityAuditor initializes the global auditor
func InitSecurityAuditor() *SecurityAuditor {
	if globalAuditor == nil {
		globalAuditor = NewSecurityAuditor()
	}
	return globalAuditor
}

// GetAuditor returns the global auditor instance
func GetAuditor() *SecurityAuditor {
	if globalAuditor == nil {
		globalAuditor = InitSecurityAuditor()
	}
	return globalAuditor
}
