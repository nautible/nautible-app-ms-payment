package jp.co.ogis_ri.nautible.app.payment.convenience.domain;

public class Convenience {
    
    private String acceptNo;
    
    private String receptDate;

    private String status; // 01:未決済 02:決済済み 09:キャンセル

    public String getAcceptNo() {
        return acceptNo;
    }

    public void setAcceptNo(String acceptNo) {
        this.acceptNo = acceptNo;
    }

    public String getReceptDate() {
        return receptDate;
    }

    public void setReceptDate(String receptDate) {
        this.receptDate = receptDate;
    }

    public String getStatus() {
        return status;
    }

    public void setStatus(String status) {
        this.status = status;
    }

}
