package jp.co.ogis_ri.nautible.app.payment.credit.domain;

public interface PaymentRepository {

    Payment getByPaymentNo(String paymentNo);

    Payment create(Payment payment) throws PaymentException;

    Payment delete(String paymentId);

    Payment update(Payment payment);
}
