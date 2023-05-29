package {{.Package}};

import org.osgi.service.component.annotations.Component;

@Component(
    immediate = true,
    property = {
        "osgi.command.scope=lfr",
        "osgi.command.function=hello"
    },
    service = {{.CamelCaseName}}.class
)
public class {{.CamelCaseName}} {

  public void hello(){
    System.out.println("hello world");
  }

}