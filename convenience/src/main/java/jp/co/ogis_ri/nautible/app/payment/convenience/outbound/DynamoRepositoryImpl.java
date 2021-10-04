package jp.co.ogis_ri.nautible.app.payment.convenience.outbound;

import java.util.HashMap;
import java.util.Map;

import javax.enterprise.context.ApplicationScoped;
import javax.inject.Inject;

import jp.co.ogis_ri.nautible.app.payment.convenience.domain.Payment;
import jp.co.ogis_ri.nautible.app.payment.convenience.domain.PaymentException;
import jp.co.ogis_ri.nautible.app.payment.convenience.domain.PaymentRepository;
import software.amazon.awssdk.enhanced.dynamodb.DynamoDbEnhancedClient;
import software.amazon.awssdk.enhanced.dynamodb.DynamoDbTable;
import software.amazon.awssdk.enhanced.dynamodb.Key;
import software.amazon.awssdk.enhanced.dynamodb.TableSchema;
import software.amazon.awssdk.services.dynamodb.DynamoDbClient;
import software.amazon.awssdk.services.dynamodb.model.AttributeAction;
import software.amazon.awssdk.services.dynamodb.model.AttributeValue;
import software.amazon.awssdk.services.dynamodb.model.AttributeValueUpdate;
import software.amazon.awssdk.services.dynamodb.model.ReturnValue;
import software.amazon.awssdk.services.dynamodb.model.UpdateItemRequest;
import software.amazon.awssdk.services.dynamodb.model.UpdateItemResponse;

@ApplicationScoped
public class DynamoRepositoryImpl implements PaymentRepository {

    private static final String PAYMENT_TABLE_NAME = "Payment";
    @Inject
    DynamoDbClient dynamoDB;

    @Override
    public Payment getByPaymentNo(String paymentNo) {
        Key key = Key.builder().partitionValue(paymentNo).build();
        Payment result = getDynamoDbTable().getItem(r -> r.key(key));
        return result;
    }

    @Override
    public Payment create(Payment payment) throws PaymentException {
        int sequence = getSequenceNumber(PAYMENT_TABLE_NAME);
        payment.setPaymentNo(Integer.toString(sequence));
        getDynamoDbTable().putItem(payment);
        return payment;
    }

    @Override
    public Payment delete(String paymentId) {
        Key key = Key.builder().partitionValue(paymentId).build();
        return getDynamoDbTable().deleteItem(key);
    }

    @Override
    public Payment update(Payment payment) {
        getDynamoDbTable().updateItem(payment);
        return payment;
    }

    private int getSequenceNumber(String tableName) {
        Map<String, AttributeValue> key = new HashMap<>();
        key.put("Name", AttributeValue.builder().s(tableName).build());
        Map<String, AttributeValueUpdate> update = new HashMap<>();
        update.put("SequenceNumber", AttributeValueUpdate.builder().value(AttributeValue.builder().n("1").build())
                .action(AttributeAction.ADD).build());
        UpdateItemRequest updateRequest = UpdateItemRequest.builder().tableName("Sequence").key(key)
                .attributeUpdates(update).returnValues(ReturnValue.UPDATED_NEW).build();
        UpdateItemResponse updateResponse = dynamoDB.updateItem(updateRequest);
        return Integer.parseInt(updateResponse.attributes().get("SequenceNumber").n());
    }

    private DynamoDbTable<Payment> getDynamoDbTable() {
        DynamoDbEnhancedClient enhancedClient = DynamoDbEnhancedClient.builder().dynamoDbClient(dynamoDB).build();
        DynamoDbTable<Payment> mappedTable = enhancedClient.table(PAYMENT_TABLE_NAME,
                TableSchema.fromBean(Payment.class));
        return mappedTable;
    }
}
