package jp.co.ogis_ri.nautible.app.payment.cash.domain;

import software.amazon.awssdk.enhanced.dynamodb.mapper.annotations.DynamoDbAttribute;
import software.amazon.awssdk.enhanced.dynamodb.mapper.annotations.DynamoDbBean;
import software.amazon.awssdk.enhanced.dynamodb.mapper.annotations.DynamoDbPartitionKey;

@DynamoDbBean
public class Payment {
    private String paymentNo = null;
    private String orderNo = null;
    private String orderDate = null;
    private Integer customerId = null;
    private Integer totalPrice = null;
    private String orderStatus = null;

    public void setOrderNo(String orderNo) {
        this.orderNo = orderNo;
    }

    public void setOrderDate(String orderDate) {
        this.orderDate = orderDate;
    }

    public void setCustomerId(Integer customerId) {
        this.customerId = customerId;
    }

    public void setTotalPrice(Integer totalPrice) {
        this.totalPrice = totalPrice;
    }

    public void setOrderStatus(String orderStatus) {
        this.orderStatus = orderStatus;
    }
    
    @DynamoDbPartitionKey
    @DynamoDbAttribute("PaymentNo")
    public String getPaymentNo() {
        return paymentNo;
    }

    public void setPaymentNo(String paymentNo) {
        if (this.paymentNo != null) {
            throw new PaymentException("W0001", "this payment object is already in use. ", null);
        }
        this.paymentNo = paymentNo;
    }

    @DynamoDbAttribute("OrderNo")
    public String getOrderNo() {
        return orderNo;
    }

    @DynamoDbAttribute("OrderDate")
    public String getOrderDate() {
        return orderDate;
    }

    @DynamoDbAttribute("CustomerId")
    public Integer getCustomerId() {
        return customerId;
    }

    @DynamoDbAttribute("TotalPrice")
    public Integer getTotalPrice() {
        return totalPrice;
    }

    @DynamoDbAttribute("OrderStatus")
    public String getOrderStatus() {
        return orderStatus;
    }
}
