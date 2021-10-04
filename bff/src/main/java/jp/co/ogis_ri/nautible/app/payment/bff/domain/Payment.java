package jp.co.ogis_ri.nautible.app.payment.bff.domain;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

/**
 * 決済データを各バックエンド（代引き、コンビニ、クレジット）とやり取りするためのオブジェクト
 * 不要な項目は無視するようにJsonIgnorePropertiesを付与
 */
@JsonIgnoreProperties(ignoreUnknown = true)
public class Payment {
    private String paymentNo = null;
    private String orderNo = null;
    private String orderDate = null;
    private Integer customerId = null;
    private Integer totalPrice = null;
    private String orderStatus = null;
    private String acceptNo = null;
    private String receiptDate = null;
    private String paymentType = null;
    private String requestId = null;
    
    public String getPaymentNo() {
        return paymentNo;
    }
    public void setPaymentNo(String paymentNo) {
        this.paymentNo = paymentNo;
    }
    public String getOrderNo() {
        return orderNo;
    }
    public void setOrderNo(String orderNo) {
        this.orderNo = orderNo;
    }
    public String getOrderDate() {
        return orderDate;
    }
    public void setOrderDate(String orderDate) {
        this.orderDate = orderDate;
    }
    public Integer getCustomerId() {
        return customerId;
    }
    public void setCustomerId(Integer customerId) {
        this.customerId = customerId;
    }
    public Integer getTotalPrice() {
        return totalPrice;
    }
    public void setTotalPrice(Integer totalPrice) {
        this.totalPrice = totalPrice;
    }
    public String getOrderStatus() {
        return orderStatus;
    }
    public void setOrderStatus(String orderStatus) {
        this.orderStatus = orderStatus;
    }
    public String getAcceptNo() {
        return acceptNo;
    }
    public void setAcceptNo(String acceptNo) {
        this.acceptNo = acceptNo;
    }
    public String getReceiptDate() {
        return receiptDate;
    }
    public void setReceiptDate(String receiptDate) {
        this.receiptDate = receiptDate;
    }
    public String getRequestId() {
        return requestId;
    }
    public void setRequestId(String requestId) {
        this.requestId = requestId;
    }
    public String getPaymentType() {
        return paymentType;
    }
    public void setPaymentType(String paymentType) {
        this.paymentType = paymentType;
    }

}
