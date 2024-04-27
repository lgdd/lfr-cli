{{- if eq .TemplateEngine "thymeleaf" -}}
<?xml version="1.0" encoding="UTF-8"?>
<div xmlns:th="http://www.thymeleaf.org">
	<p class="user-greeting" th:text="#{greetings(${user.firstName},${user.lastName})}" />
	<p th:text="#{todays-date-is(${todaysDate})}" />
</div>
{{- else  -}}
<?xml version="1.0" encoding="UTF-8"?>
<jsp:root xmlns:jsp="http://java.sun.com/JSP/Page"
          xmlns:portlet="http://xmlns.jcp.org/portlet_3_0"
          xmlns:spring="http://www.springframework.org/tags"
          version="2.1">
    <jsp:directive.page contentType="text/html" pageEncoding="UTF-8" />
    <portlet:defineObjects/>
    <p class="user-greeting">
        <spring:message arguments="${user.firstName},${user.lastName}" code="greetings" />
    </p>
    <p>
        <spring:message arguments="${todaysDate}" argumentSeparator=";" code="todays-date-is" />
    </p>
</jsp:root>
{{- end  -}}
