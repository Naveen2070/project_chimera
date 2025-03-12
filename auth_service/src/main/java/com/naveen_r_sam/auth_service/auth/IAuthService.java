package com.naveen_r_sam.auth_service.auth;

import com.naveen_r_sam.auth_service.model.Users;
import org.springframework.http.ResponseEntity;

public interface IAuthService {
    /**
     * Register a user, given a {@link Users} object.
     * @param user the user to register
     * @return a response entity containing a JSON object with
     * a single key, "message", containing a message describing the
     * result of the registration
     */
    ResponseEntity<?> registerUser(Users user);
    /**
     * Authenticate a user, given a {@link Users} object.
     * @param user the user to authenticate
     * @return a response entity containing a JSON object with
     * a single key, "message", describing the result of the authentication,
     * and possibly additional authentication information if successful
     */
    ResponseEntity<?> authenticateUser(Users user);
}
