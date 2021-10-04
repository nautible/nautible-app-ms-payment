package jp.co.ogis_ri.nautible.app.payment.bff.outbound.dynamodb;

import java.util.List;

import javax.enterprise.context.ApplicationScoped;

import jp.co.ogis_ri.nautible.app.payment.bff.domain.Payment;
import jp.co.ogis_ri.nautible.app.payment.bff.domain.PaymentRepository;

@ApplicationScoped
public class RestPaymentRepositoryImpl implements PaymentRepository {

    @Override
    public Payment getByPaymentNo(String paymentNo) {
        // TODO Auto-generated method stub
        return null;
    }

    @Override
    public List<Payment> getByCustomerIdAndTerm(Integer customerId, String orderDateFrom, String orderDateTo) {
        // TODO Auto-generated method stub
        return null;
    }

    @Override
    public Payment getByRequestId(String requestId) {
        // TODO Auto-generated method stub
        return null;
    }

}
