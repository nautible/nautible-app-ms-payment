package jp.co.ogis_ri.nautible.app.payment.credit.inbound.rest;

import java.util.List;

import org.mapstruct.Mapper;
import org.mapstruct.factory.Mappers;

import jp.co.ogis_ri.nautible.app.payment.credit.config.QuarkusMappingConfig;
import jp.co.ogis_ri.nautible.app.payment.credit.domain.Payment;

@Mapper(config = QuarkusMappingConfig.class)
public interface PaymentMapper {

    PaymentMapper INSTANCE = Mappers.getMapper(PaymentMapper.class);

    RestPayment paymentToRestPayment(Payment payment);

    List<RestPayment> paymentToRestPayment(List<Payment> payment);

    Payment restPaymentToPayment(RestPayment payment);

    Payment restUpdatePaymentToPayment(RestUpdatePayment payment);

    Payment restCreatePaymentToPayment(RestCreatePayment payment);

    List<Payment> restPaymentToPayment(List<Payment> payment);

}
