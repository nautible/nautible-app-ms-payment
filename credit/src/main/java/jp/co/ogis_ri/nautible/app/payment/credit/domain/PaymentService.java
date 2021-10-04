package jp.co.ogis_ri.nautible.app.payment.credit.domain;

import javax.enterprise.context.ApplicationScoped;
import javax.inject.Inject;

import com.google.gson.Gson;

@ApplicationScoped
public class PaymentService {
    @Inject
    PaymentRepository paymentRepository;

    @Inject
    CreditRepository creditRepository;

    public Payment get(String paymentNo) {
        return paymentRepository.getByPaymentNo(paymentNo);
    }

    public Payment create(Payment payment) throws PaymentException {
        // MDK部分のモック
        String response = creditRepository.create(payment.getOrderNo(), payment.getOrderDate(), payment.getCustomerId(), payment.getTotalPrice());
        Gson gson = new Gson();
        Credit credit = gson.fromJson(response, Credit.class);
        payment.setOrderStatus(credit.getStatus());
        payment.setAcceptNo(credit.getAcceptNo());
        payment.setReceiptDate(credit.getReceptDate());
        return paymentRepository.create(payment);
    }

    public Payment update(Payment payment) {
        if (paymentRepository.getByPaymentNo(payment.getPaymentNo()) == null) {
            return null;
        }
        return paymentRepository.update(payment);
    }

    public Payment deleteByPaymentNo(String paymentNo) {
        Payment payment = paymentRepository.getByPaymentNo(paymentNo);
        if (payment == null) {
            return null;
        }
        // MDK部分のモック
        String response = creditRepository.cancel(payment.getAcceptNo(), payment.getReceiptDate());
        Gson gson = new Gson();
        Credit credit = gson.fromJson(response, Credit.class);
        payment.setOrderStatus(credit.getStatus());

        return paymentRepository.delete(paymentNo);
    }
}
