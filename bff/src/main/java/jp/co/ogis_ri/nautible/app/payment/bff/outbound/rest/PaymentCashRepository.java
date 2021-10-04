package jp.co.ogis_ri.nautible.app.payment.bff.outbound.rest;

import javax.ws.rs.DELETE;
import javax.ws.rs.POST;
import javax.ws.rs.PUT;
import javax.ws.rs.Path;

import org.eclipse.microprofile.rest.client.inject.RegisterRestClient;
import org.jboss.resteasy.annotations.jaxrs.PathParam;

import jp.co.ogis_ri.nautible.app.payment.bff.domain.Payment;

/**
 * 決済処理をバックエンド（代引き）とやり取りするインターフェース
 */
@Path("/nautible-app-payment-cash")
@RegisterRestClient
public interface PaymentCashRepository {

    @POST
    @Path("/method/payment/")
    Payment create(Payment payment);

    @DELETE
    @Path("/method/payment/{paymentNo}")
    Payment delete(@PathParam String paymentNo);

    @PUT
    @Path("/methodpayment/")
    Payment update(Payment payment);

    //Payment getByRequestId(String requestId);
}
