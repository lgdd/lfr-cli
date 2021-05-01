{{- if  eq .TemplateEngine "thymeleaf" }}
<?xml version="1.0" encoding="UTF-8"?>
<div xmlns:th="http://www.thymeleaf.org">
	<form th:id="|${namespace}mainForm|" th:action="${mainFormActionURL}" class="user-form" method="post" th:object="${user}">
		<span class="portlet-msg-error" th:if="${#fields.hasGlobalErrors()}" th:errors="*{global}"/>
		<p class="caption" th:text="#{personal-information}"/>
		<fieldset>
			<div class="form-group">
				<label th:for="|${namespace}firstName|" path="firstName" th:text="#{first-name}"/>
				<input th:id="|${namespace}firstName|" class="form-control" th:field="*{firstName}"/>
				<span class="portlet-msg-error" th:if="${#fields.hasErrors('firstName')}" th:errors="*{firstName}"/>
			</div>
			<div class="form-group">
				<label th:for="|${namespace}lastName|" path="lastName" th:text="#{last-name}"/>
				<input th:id="|${namespace}lastName|" class="form-control" th:field="*{lastName}"/>
				<span class="portlet-msg-error" th:if="${#fields.hasErrors('lastName')}" th:errors="*{lastName}"/>
			</div>
		</fieldset>
		<hr />
		<input class="btn btn-primary" th:text="${submit}" type="submit"/>
	</form>
</div>
{{- else }}
<?xml version="1.0" encoding="UTF-8"?>
<jsp:root xmlns:jsp="http://java.sun.com/JSP/Page"
          xmlns:portlet="http://xmlns.jcp.org/portlet_3_0"
          xmlns:spring="http://www.springframework.org/tags"
          xmlns:form="http://www.springframework.org/tags/form"
          version="2.1">
    <jsp:directive.page contentType="text/html" pageEncoding="UTF-8" />
    <portlet:defineObjects/>
    <portlet:actionURL var="mainFormActionURL"/>
    <form:form id="${namespace}mainForm" action="${mainFormActionURL}" class="user-form" method="post" modelAttribute="user">
        <form:errors cssClass="portlet-msg-error" />
        <p class="caption">
            <spring:message code="personal-information" />
        </p>
        <fieldset>
            <div class="form-group">
                <form:label for="${namespace}firstName" path="firstName">
                    <spring:message code="first-name" />
                </form:label>
                <form:input id="${namespace}firstName" cssClass="form-control" path="firstName"/>
                <form:errors path="firstName" cssClass="portlet-msg-error"/>
            </div>
            <div class="form-group">
                <form:label for="${namespace}lastName" path="lastName">
                    <spring:message code="last-name" />
                </form:label>
                <form:input id="${namespace}lastName" cssClass="form-control" path="lastName"/>
                <form:errors path="lastName" cssClass="portlet-msg-error"/>
            </div>
        </fieldset>
        <hr />
        <spring:message code="submit" var="submit" />
        <input class="btn btn-primary" value="${submit}" type="submit"/>
    </form:form>
</jsp:root>
{{- end  }}
