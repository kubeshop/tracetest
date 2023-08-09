# Expressions

Tracetest allows you to add expressions when writing your tests. They are a nice and clean way of adding values that are only known during execution time. For example, when referencing a variable, a span attribute or even arithmetic operations.

## **Features**

* Reference span attributes
* Reference variables
* Arithmetic operations
* String interpolation
* Filters

### **Reference Span Attributes**

When building assertions, you might need to assert if a certain span contains an attribute and that this attribute has a specific value. To accomplish this with Tracetest, you can use expressions to get the value of the span. When referencing an attribute, add the prefix `attr:` and its name. For example, imagine you have to check if the attribute `service.name` is equal to `cart-api`. Use the following statement:

```css
attr:service.name = "cart-api"
```

### **Reference Variables**

Create variables in Tracetest based on the trace obtained by the test to enable assertions that require values from other spans. Variables use the prefix `var:` and its name. For example, a variable called `user_id` would be referenced as `var:user_id` in an expression.

### **Arithmetic Operations**

Sometimes we need to manipulate data to ensure our test data is correct. As an example, we will use a purchase operation. How you would make sure that, after the purchase, the product inventory is smaller than before? For this, we might want to use arithmetic operations:

```css
attr:product.stock = attr:product.stok_before_purchase - attr:product.number_bought_items
```

### **String Interpolation**

Some tests might require strings to be compared, but maybe you need to generate a dynamic string that relies on a dynamic value. This might be used in an assertion or even in the request body referencing a variable.

```css
attr:error.message = "Could not withdraw ${attr:withdraw.amount}, your balance is insufficient."
```

Note that within `${}` you can add any expression, including arithmetic operations and filters.


### **Filters**

Filters are functions that are executed using the value obtained by the expression. They are useful to transform the data. Multiple filters can be chained together. The output of the previous filter will be used as the input to the next until all filters are executed.

#### **JSON Path**
This filter allows you to filter a JSON string and obtain only data that is relevant.

```css
'{ "name": "Jorge", "age": 27, "email": "jorge@company.com" }' | json_path '.age' = 27
```

If multiple values are matched, the output will be a flat array containing all values.

```css
'{ "array": [{"name": "Jorge", "age": 27}, {"name": "Tim", "age": 52}]}' | json_path '$.array[*]..["name", "age"] = '["Jorge", 27, "Tim", 52]'
```

#### **RegEx**
Filters part of the input that match a RegEx. Imagine you have a specific part of a text that you want to extract:

```css
'My account balance is $48.52' | regex '\$\d+(\.\d+)?' = '$48.52'
```

#### **RegEx Group**
If matching more than one value is required, you can define groups for your RegEx and extract multiple values at once.

Wrap the groups you want to extract with parentheses.

```css
'Hello Marcus, today you have 8 meetings' | regex_group 'Hello (\w+), today you have (\d+) meetings' = '["Marcus", 8]'
```

### **Get Index**

Some filters might result in an array. If you want to assert just part of this array, this filter allows you to pick one element from the array based on its index.

```css
'{ "array": [1, 2, 3] }' | json_path '$.array[*]' | get_index 1 = 2
```

You can select the last item from a list by specifying `'last'` as the argument for the filter:

```css
'{ "array": [1, 2, 3] }' | json_path '$.array[*]' | get_index 'last' = 3
```

### **Length**

Returns the size of the input array. If it's a single value, it will return 1. Otherwise it will return `length(input_array)`.

```css
'{ "array": [1, 2, 3] }' | json_path '$.array[*]' | length = 3
```

```css
"my string" | length = 1
```

### **Type**

Return the type of the input as a string.

```css
2 | type = "number"
```

```css
2s | type = "duration"
```

```css
"something" | type = "string"
```

```css
[1, 2s, "this is a string"] | type = "array"
```

```css
attr:myapp.operations | get_index 2 | type = "string"
```

```css
# check if attribute is either a number of a string
["number", "string"] contains attr:my_attribute | type
```
