package jp.co.ogis_ri.nautible.app.payment.credit.domain;

public interface CreditRepository {

    String create(String orderNo, String orderDate, Integer custmerId, Integer totalPrice);

    String cancel(String acceptNo, String receptDate);

    String update(String acceptNo, String receptDate, String orderNo, String orderDate, String custmerId, Integer totalPrice);
}
