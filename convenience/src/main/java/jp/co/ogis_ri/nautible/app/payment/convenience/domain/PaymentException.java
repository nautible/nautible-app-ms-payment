package jp.co.ogis_ri.nautible.app.payment.convenience.domain;

public class PaymentException extends RuntimeException {
    
    private String errorCode;

    public String getErrorCode() {
        return errorCode;
    }

    public PaymentException(String errorCode, String message, Throwable cause) {
        super(message, cause);
        this.errorCode = errorCode;
    }
}
