package jp.co.ogis_ri.nautible.app.payment.cash.domain;

import javax.enterprise.context.ApplicationScoped;
import javax.inject.Inject;

@ApplicationScoped
public class PaymentService {
    @Inject
    PaymentRepository paymentRepository;

    public Payment get(String paymentNo) {
        return paymentRepository.getByPaymentNo(paymentNo);
    }

    public Payment create(Payment payment) throws PaymentException {
        payment.setOrderStatus("01");
        return paymentRepository.create(payment);
    }

    public Payment update(Payment payment) {
        if (paymentRepository.getByPaymentNo(payment.getPaymentNo()) == null) {
            return null;
        }
        return paymentRepository.update(payment);
    }

    public Payment deleteByPaymentNo(String paymentNo) {
        return paymentRepository.delete(paymentNo);
    }
}
