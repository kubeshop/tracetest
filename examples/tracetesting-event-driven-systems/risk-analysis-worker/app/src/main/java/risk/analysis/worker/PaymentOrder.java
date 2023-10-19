package risk.analysis.worker;

public class PaymentOrder {
  private String originCustomerID;
  private String destinationCustomerID;
  private Double value;

  public String getOriginCustomerID() {
    return originCustomerID;
  }
  public void setOriginCustomerID(String originCustomerID) {
    this.originCustomerID = originCustomerID;
  }
  public String getDestinationCustomerID() {
    return destinationCustomerID;
  }
  public void setDestinationCustomerID(String destinationCustomerID) {
    this.destinationCustomerID = destinationCustomerID;
  }
  public Double getValue() {
    return value;
  }
  public void setValue(Double value) {
    this.value = value;
  }
}
