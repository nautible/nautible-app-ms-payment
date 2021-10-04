package jp.co.ogis_ri.nautible.app.payment.bff.outbound.rest;

import javax.ws.rs.DELETE;
import javax.ws.rs.POST;
import javax.ws.rs.PUT;
import javax.ws.rs.Path;

import org.eclipse.microprofile.rest.client.inject.RegisterRestClient;
import org.jboss.resteasy.annotations.jaxrs.PathParam;

import jp.co.ogis_ri.nautible.app.payment.bff.domain.Payment;

/**
 * 決済処理をバックエンド（クレジット）とやり取りするためのインターフェース
 */
@Path("/nautible-app-payment-credit")
@RegisterRestClient
public interface PaymentCreditRepository {

    @POST
    @Path("/method")
    Payment create(Payment payment);

    @DELETE
    @Path("/method/{paymentNo}")
    Payment delete(@PathParam String paymentNo);

    @PUT
    @Path("/method")
    Payment update(Payment payment);

    //Payment getByRequestId(String requestId);
}
