// Package scaffold generates the file and directory structure for Liferay
// workspace modules (API, MVC portlet, Spring portlet, REST Builder, Service
// Builder, Gogo shell command), client extensions, and Docker configuration
// files. Each Create* function copies embedded templates, renders them with
// Go's text/template engine, and reports the created files via the logger.
package scaffold
