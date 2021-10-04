package jp.co.ogis_ri.nautible.app.payment.bff.domain;

import java.util.function.Function;

import javax.enterprise.context.ApplicationScoped;

import io.dapr.client.DaprClient;
import io.dapr.client.DaprClientBuilder;
import io.dapr.client.domain.PublishEventRequestBuilder;
import jp.co.ogis_ri.nautible.app.order.inbound.rest.CreateOrderReply;
import jp.co.ogis_ri.nautible.app.order.inbound.rest.CreateOrderReply.ProcessTypeEnum;

/**
 * 決済のSAGAマネージャ
 */
@ApplicationScoped
public class PaymentSagaManager {
    /** orderのpubsub名 */
    private static final String ORDER_PUBSUB_NAME = "order-pubsub";
    // デフォルトのapplication/jsonだとsubscriber側でCloudEventに変換するとデータ部がMapになる。
    // steamにしてObjectMapperでマッピングした方が扱いやすい。
    private static final String PUBSUB_CONTENT_TYPE = "application/octet-stream";

    /**
     * 決済処理の正常応答を返す
     * @param requestId リクエストID
     */
    public void replyCreate(String requestId) {
        reply(requestId);
    }

    /**
     * 決済処理のリクエスト不正の応答を返す
     * @param requestId リクエストID
     * @param message メッセージ
     */
    public void replyCreateBadRequest(String requestId, String message) {
        replyBadRequest(requestId, message);
    }

    /**
     * 決済取消し処理の正常応答を返す
     * @param requestId リクエストID
     */
    public void replyRejectCreate(String requestId) {
        reply(requestId);
    }

    /**
     * 決済処理のリクエスト不正の応答を返す
     * @param requestId リクエストID
     * @param message メッセージ
     */
    public void replyRejectCreateBadRequest(String requestId, String message) {
        replyBadRequest(requestId, message);
    }

    /**
     * 正常応答を返す
     * @param requestId リクエストID
     */
    private void reply(String requestId) {
        reply(requestId, "create-order-reply",
                new CreateOrderReply().status(200).requestId(requestId)
                        .processType(ProcessTypeEnum.PAYMENT));
    }

    /**
     * リクエスト不正の応答を返す
     * @param requestId リクエストID
     * @param message メッセージ
     */
    private void replyBadRequest(String requestId, String message) {
        reply(requestId, "create-order-reply",
                new CreateOrderReply().status(400).requestId(requestId)
                        .processType(ProcessTypeEnum.PAYMENT).message(message));
    }

    /**
     * SAGA。Daprのpubsubを利用して、応答のpublishを行う。
     * @param requestId リクエストId
     * @param replyTopic 応答のtopic
     * @param data 応答data
     */
    private void reply(String requestId, String replyTopic, Object data) {
        // WARNING 分散トレーシングのIstio/Daprが共存できない。個別には実現できる。
        // DaprはW3Cのspecを採用、IstioはW3CのSpecには現状未対応。
        // IstioがW3Cに対応したらうまく共存できるかも？https://github.com/istio/istio/issues/23960
        executeDaprClient(
                c -> c.publishEvent(new PublishEventRequestBuilder(ORDER_PUBSUB_NAME, replyTopic, data)
                        .withContentType(PUBSUB_CONTENT_TYPE).build())
                        .block());
    }

    private <R> R executeDaprClient(Function<DaprClient, R> func) {
        try (DaprClient client = new DaprClientBuilder().build()) {
            return func.apply(client);
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

}
