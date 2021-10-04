package jp.co.ogis_ri.nautible.app.payment.bff.domain;

import java.util.List;

/**
 * 決済データをDBから取得するインターフェース
 */
public interface PaymentRepository {

    Payment getByPaymentNo(String paymentNo);

    List<Payment> getByCustomerIdAndTerm(Integer customerId, String orderDateFrom, String orderDateTo);
    
    Payment getByRequestId(String requestId);

}
