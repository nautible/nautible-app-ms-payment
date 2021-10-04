package jp.co.ogis_ri.nautible.app.payment.bff.inbound.rest;

import java.io.IOException;
import java.util.List;
import java.util.Set;
import java.util.logging.Level;
import java.util.logging.Logger;

import javax.inject.Inject;
import javax.validation.ConstraintViolation;
import javax.validation.Valid;
import javax.validation.Validator;
import javax.ws.rs.core.Response;
import javax.ws.rs.core.Response.Status;

import com.fasterxml.jackson.databind.ObjectMapper;

import io.dapr.client.domain.CloudEvent;
import jp.co.ogis_ri.nautible.app.payment.bff.domain.Payment;
import jp.co.ogis_ri.nautible.app.payment.bff.domain.PaymentService;

public class RestPaymentServiceImpl implements RestPaymentService {

    Logger LOG = Logger.getLogger(RestPaymentServiceImpl.class.getName());
    /** paymentのpubsub */
    private static final String PAYMENT_PUBSUB_NAME = "payment-pubsub";

    @Inject
    PaymentMapper mapper;
    @Inject
    ObjectMapper objectMapper;
    @Inject
    PaymentService service;
    @Inject
    Validator validator;

    @Override
    public Response delete(String paymentNo) {
        Payment payment = service.getByPaymentNo(paymentNo);
        Payment paymentRet = service.deleteByPaymentNo(payment.getPaymentType(), paymentNo);
        return paymentRet == null ? Response.status(Status.NOT_FOUND).build()
                : Response.status(Status.OK).build();
    }

    @Override
    public Response find(Integer customerId, String orderDateFrom, String orderDateTo) {
        List<Payment> payments = service.getByCustomerIdAndTerm(customerId, orderDateFrom, orderDateTo);
        return (payments == null || payments.size() == 0) ? Response.status(Status.NOT_FOUND).build()
                : Response.ok(mapper.paymentToRestPayment(payments)).build();
    }

    @Override
    public Response getByPaymentNo(String paymentNo) {
        Payment payment = service.getByPaymentNo(paymentNo);
        return payment == null ? Response.status(Status.NOT_FOUND).build()
                : Response.ok(mapper.paymentToRestPayment(payment)).build();
    }

    @Override
    public Response update(@Valid RestUpdatePayment restUpdatePayment) {
        Payment paymentRet = service.update(mapper.restUpdatePaymentToPayment(restUpdatePayment));
        return paymentRet == null ? Response.status(Status.NOT_FOUND).build()
                : Response.ok(mapper.paymentToRestPayment(paymentRet)).build();
    }

    @Override
    public Response daprSubscribe() {
        // https://docs.dapr.io/developing-applications/building-blocks/pubsub/howto-publish-subscribe/#programmatic-subscriptions
        return Response.ok(List.of(
                new DaprSubscribe().pubsubname(PAYMENT_PUBSUB_NAME).topic("payment-create")
                        .route("/payment/create"),
                new DaprSubscribe().pubsubname(PAYMENT_PUBSUB_NAME)
                        .topic("payment-reject-create").route("/payment/rejectCreate")))
                .build();
    }

    @Override
    public Response create(@Valid byte[] body) {
        final RestCreatePayment request = readCloudEventRequest(body, RestCreatePayment.class);
        LOG.log(Level.INFO, request.toString());
        Set<ConstraintViolation<RestCreatePayment>> violations = validator.validate(request);
        if (violations.isEmpty() == false) {
            service.replyCreateBadRequest(request.getRequestId(), "Bad request.");
            return Response.ok().build();
        }
        service.create(mapper.restCreatePaymentToPayment(request));

        return Response.ok().build();
    }

    @Override
    public Response rejectCreate(@Valid byte[] body) {
        final RestRejectCreatePayment request = readCloudEventRequest(body, RestRejectCreatePayment.class);
        Set<ConstraintViolation<RestRejectCreatePayment>> violations = validator.validate(request);
        if (violations.isEmpty() == false) {
            service.replyRejectCreateBadRequest(request.getRequestId(), "Bad request.");
            return Response.ok().build();
        }
        service.rejectCreate(request.getRequestId());

        return Response.ok().build();
    }

    private <R> R readCloudEventRequest(byte[] body, Class<R> clazz) {
        try {
            CloudEvent event = CloudEvent.deserialize(body);
            return objectMapper.readValue(event.getBinaryData(), clazz);
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
    }

}