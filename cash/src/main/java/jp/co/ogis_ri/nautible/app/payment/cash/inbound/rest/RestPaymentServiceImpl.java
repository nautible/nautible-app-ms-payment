package jp.co.ogis_ri.nautible.app.payment.cash.inbound.rest;

import java.util.logging.Logger;

import javax.inject.Inject;
import javax.validation.Valid;
import javax.validation.Validator;
import javax.ws.rs.core.Response;
import javax.ws.rs.core.Response.Status;

import com.fasterxml.jackson.databind.ObjectMapper;

import jp.co.ogis_ri.nautible.app.payment.cash.domain.Payment;
import jp.co.ogis_ri.nautible.app.payment.cash.domain.PaymentException;
import jp.co.ogis_ri.nautible.app.payment.cash.domain.PaymentService;

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
        Payment payment = service.get(paymentNo);
        return payment == null ? Response.status(Status.NOT_FOUND).build()
                : Response.ok(mapper.paymentToRestPayment(payment)).build();
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
        Payment payment = service.update(mapper.restUpdatePaymentToPayment(request));
        return payment == null ? Response.status(Status.NOT_FOUND).build()
                : Response.ok(mapper.paymentToRestPayment(payment)).build();
    }

    @Override
    public Response delete(String paymentNo) {
        Payment payment = service.deleteByPaymentNo(paymentNo);
        return payment == null ? Response.status(Status.NOT_FOUND).build()
                : Response.status(Status.OK).build();
    }
}