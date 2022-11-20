# Versioning

As your system evolves, your tests tend to do the same. However, that might be confusing if you don't have a versioning mechanism in place. Imagine that you wrote a new test for the version `v0.5.0` of your application. After some months, your application is in version `v0.13.7`. Most likely, your tests changed as you moved your application forward. But without versioning, if you revisit that first test you created, it will look like exactly the one you use today instead of the test you originally wrote. That happens because while you have multiple versions of your application, you only keep track of one version of your tests: the current version. So there is no way of going back in time and seeing what a test looked like in the past.

**But that is not a problem if you use Tracetest. It has versioning built-in!**

## **How It Works**
Once you create a test, it is tagged as the initial version (`v1`). Every time you change something in your test (edit its identification details, add assertions, change selectors, etc) Tracetest detects those changes and increase the version by 1. If no changes were made, the version is kept untouched.

### **Change Detection**
These are the fields of a test that are checked to verify if it has changed:

* name
* description
* trigger
* test definition
    * selectors
    * assertions

### **Problems**
If you notice you are editing fields and your test version is not changing, let us know by opening a [bug report](https://github.com/kubeshop/tracetest/issues) on our Github Repository.