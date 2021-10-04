package jp.co.ogis_ri.nautible.app.payment.convenience.outbound;

import java.text.SimpleDateFormat;
import java.util.Calendar;
import java.util.UUID;

import javax.enterprise.context.ApplicationScoped;

import com.google.gson.Gson;

import jp.co.ogis_ri.nautible.app.payment.convenience.domain.ConvenienceRepository;

@ApplicationScoped
public class MockConvenienceRepositoryImpl implements ConvenienceRepository {


    @Override
    public String create(String orderNo, String orderDate, Integer custmerId, Integer totalPrice) {
        SimpleDateFormat sdf = new SimpleDateFormat("yyyy/MM/dd");
        CreditResponse response = new CreditResponse(UUID.randomUUID().toString(), sdf.format(Calendar.getInstance().getTime()), "01");
        Gson gson = new Gson();
        return gson.toJson(response);
    }

    @Override
    public String cancel(String acceptNo, String receptDate) {
        CreditResponse response = new CreditResponse(acceptNo, receptDate, "09");
        Gson gson = new Gson();
        return gson.toJson(response);
    }

    @Override
    public String update(String acceptNo, String receptDate, String orderNo, String orderDate, String custmerId, Integer totalPrice) {
        CreditResponse response = new CreditResponse(acceptNo, receptDate, "01");
        Gson gson = new Gson();
        return gson.toJson(response);
    }

    public class CreditResponse {

        private String acceptNo;
    
        private String receptDate;
    
        private String status; // 01:未決済 02:決済済み 09:キャンセル

        public CreditResponse(String acceptNo, String receptDate, String status) {
            this.acceptNo = acceptNo;
            this.receptDate = receptDate;
            this.status = status;
        }
    
        public String getAcceptNo() {
            return this.acceptNo;
        }
    
        public String getReceptDate() {
            return this.receptDate;
        }

        public String getStatus() {
            return this.status;
        }

    }
}
