# Assertions

Test Specifications may be added to a trace to set a value for a step in the trace to determine success or failure. If test specs have already been added to a test, they will be on the Test screen:

![Test Spec List](../img/test-spec-list-0.11.png)

After you have created a test and your test run is complete, click the **Add Test Spec** button at the bottom right of the Test screen.

![Add Test Spec](../img/add-test-spec-0.11.png)

The **Add Test Spec** dialog opens.

![Create Test Spec](../img/create-test-spec-0.11.png)

The span that the new test spec will apply to is hightlighted in the graph view on the left:

![Selected Span](../img/selected-span-0.11.png)

To add an assertion to a span, click the first drop down to see the list of attributes that apply to the selected span:

![Assertion Attributes](../img/assertion-attributes-0.11.png)

Then select the operator for your assertion:

![Assertion Operators](../img/assertion-operators-0.11.png)

And add the value for comparison:

![Assertion Values](../img/assertion-values-0.11.png)

Finally, you can give your test spec an optional name and click **Save Test Spec**:

![Save Test Spec](../img/save-test-spec-0.11.png)


<!--- You can also create assertions by hovering over the **+** sign to the right of an attribute in the trace. 

![Add Assertion Hover](../img/add-assertion-hover-0.6.png)

This will populate the assertion with the correct information for that attribute.

![Add Assertion Hover Details](../img/add-assertion-hover-details-0.6.png)

The **Filter** field allows for limiting the spans affected by the assertion.

![Filter Assertion](../img/assertion-filter-0.6.png)

Use the **Advanced mode** toggle switch to use the wizard or the query language to create the span selector:

![Span Selector Advanced Mode](../img/span-advanced-mode-0.6.png)

![Span Selector Advanced Mode](../img/span-advanced-mode-0.6.gif)

<!--- To see adding assertions in action, please watch <Add link to video> --->
