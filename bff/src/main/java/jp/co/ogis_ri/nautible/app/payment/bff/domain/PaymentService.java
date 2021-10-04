package jp.co.ogis_ri.nautible.app.payment.bff.domain;

import java.util.List;
import java.util.logging.Level;
import java.util.logging.Logger;

import javax.enterprise.context.ApplicationScoped;
import javax.inject.Inject;

import org.eclipse.microprofile.rest.client.inject.RestClient;

import jp.co.ogis_ri.nautible.app.payment.bff.outbound.rest.PaymentCashRepository;
import jp.co.ogis_ri.nautible.app.payment.bff.outbound.rest.PaymentConvenienceRepository;
import jp.co.ogis_ri.nautible.app.payment.bff.outbound.rest.PaymentCreditRepository;

@ApplicationScoped
public class PaymentService {

    Logger LOG = Logger.getLogger(PaymentService.class.getName());

    @Inject
    PaymentRepository paymentRepository;

    @Inject
    @RestClient
    PaymentCashRepository paymentCashRepository;

    @Inject
    @RestClient
    PaymentConvenienceRepository paymentConvenienceRepository;

    @Inject
    @RestClient
    PaymentCreditRepository paymentCreditRepository;

    @Inject
    PaymentSagaManager sagaManager;

    public Payment getByPaymentNo(String paymentNo) {
        return paymentRepository.getByPaymentNo(paymentNo);
    }

    public List<Payment> getByCustomerIdAndTerm(Integer customerId, String orderDateFrom, String orderDateTo) {
        return paymentRepository.getByCustomerIdAndTerm(customerId, orderDateFrom, orderDateTo);
    }

    public void create(Payment payment) {
        switch(payment.getPaymentType()) {
        case "01":
            LOG.log(Level.INFO, "send Cash");
            Payment paymentRet = paymentCashRepository.create(payment);
            LOG.log(Level.INFO, paymentRet.getAcceptNo());
            break;
        case "02":
            LOG.log(Level.INFO, "send Convenience");
            paymentConvenienceRepository.create(payment);
            break;
        case "03":
            LOG.log(Level.INFO, "send Credit");
            paymentCreditRepository.create(payment);
            break;
        default:
            replyCreateBadRequest(payment.getRequestId(), "決済区分が不正です");
            return;
        }
        sagaManager.replyCreate(payment.getRequestId());
    }

    public void rejectCreate(String requestId) {
        Payment payment = paymentRepository.getByRequestId(requestId);
        if (payment == null) {
            // pubsub冪等性対応
            return;
        }
        switch(payment.getPaymentType()) {
            case "01":
                paymentCashRepository.delete(payment.getPaymentNo());
                break;
            case "02":
                paymentConvenienceRepository.delete(payment.getPaymentNo());
                break;
            case "03":
                paymentCreditRepository.delete(payment.getPaymentNo());
                break;
            default:
                replyRejectCreateBadRequest(payment.getRequestId(), "決済区分が不正です");
                return;
        }
        sagaManager.replyRejectCreate(requestId);
    }

    public void replyCreateBadRequest(String requestId, String message) {
        sagaManager.replyCreateBadRequest(requestId, message);
    }

    public void replyRejectCreateBadRequest(String requestId, String message) {
        sagaManager.replyRejectCreateBadRequest(requestId, message);
    }

    public Payment update(Payment payment) {
        Payment response = null;
        if (paymentRepository.getByPaymentNo(payment.getPaymentNo()) == null) {
            return null;
        }
        switch(payment.getPaymentType()) {
            case "01":
                response = paymentCashRepository.update(payment);
                break;
            case "02":
                response = paymentConvenienceRepository.update(payment);
                break;
            case "03":
                response = paymentCreditRepository.update(payment);
                break;
            default:
                break;
        }
        return response;
    }

    public Payment deleteByPaymentNo(String paymentType, String paymentNo) {
        Payment response = null;
        switch(paymentType) {
            case "01":
                response = paymentCashRepository.delete(paymentNo);
                break;
            case "02":
                response = paymentConvenienceRepository.delete(paymentNo);
                break;
            case "03":
                response = paymentCreditRepository.delete(paymentNo);
                break;
            default:
                break;
        }
        return response;
    }

}
