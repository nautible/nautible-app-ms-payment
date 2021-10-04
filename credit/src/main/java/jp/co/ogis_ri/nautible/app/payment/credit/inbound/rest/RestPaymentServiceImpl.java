package jp.co.ogis_ri.nautible.app.payment.credit.inbound.rest;

import java.util.logging.Logger;

import javax.inject.Inject;
import javax.validation.Valid;
import javax.validation.Validator;
import javax.ws.rs.core.Response;
import javax.ws.rs.core.Response.Status;

import com.fasterxml.jackson.databind.ObjectMapper;

import jp.co.ogis_ri.nautible.app.payment.credit.domain.Payment;
import jp.co.ogis_ri.nautible.app.payment.credit.domain.PaymentException;
import jp.co.ogis_ri.nautible.app.payment.credit.domain.PaymentService;

public class RestPaymentServiceImpl implements RestPaymentService {

    Logger LOG = Logger.getLogger(RestPaymentServiceImpl.class.getName());

    @Inject
    PaymentMapper mapper;
    @Inject
    ObjectMapper objectMapper;
    @Inject
    PaymentService service;
    @Inject
    Validator validator;

    @Override
    public Response get(String paymentNo) {
        Payment paymentRet = service.get(paymentNo);
        return paymentRet == null ? Response.status(Status.NOT_FOUND).build()
                : Response.ok(mapper.paymentToRestPayment(paymentRet)).build();
    }

    @Override
    public Response create(@Valid RestCreatePayment request) {
        try {
            Payment payment = service.create(mapper.restCreatePaymentToPayment(request));
            return Response.ok(mapper.paymentToRestPayment(payment)).build();
        } catch (PaymentException pe) {
            return Response.status(Status.INTERNAL_SERVER_ERROR).build();
        }
    }

    @Override
    public Response update(@Valid RestUpdatePayment request) {
        Payment paymentRet = service.update(mapper.restUpdatePaymentToPayment(request));
        return paymentRet == null ? Response.status(Status.NOT_FOUND).build()
                : Response.ok(mapper.paymentToRestPayment(paymentRet)).build();
    }

    @Override
    public Response delete(String paymentNo) {
        Payment paymentRet = service.deleteByPaymentNo(paymentNo);
        return paymentRet == null ? Response.status(Status.NOT_FOUND).build()
                : Response.status(Status.OK).build();
    }

}