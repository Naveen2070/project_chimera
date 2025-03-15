package com.naveen_r_sam.auth_service.auth;

import com.naveen_r_sam.auth_service.dto.LoginRequestDTO;
import com.naveen_r_sam.auth_service.dto.SignUpDTO;
import com.naveen_r_sam.auth_service.model.Users;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

@RestController
public class AuthController {

    private final IAuthService authService;

    public AuthController(IAuthService authService) {
        this.authService = authService;
    }

    @PostMapping("/register")
    @ResponseBody
    public ResponseEntity<?> registerUser(@RequestBody SignUpDTO user) {
        return authService.registerUser(user);
    }

    @PostMapping("/login")
    public ResponseEntity<?> login(@RequestBody LoginRequestDTO user) {
        return authService.authenticateUser(user);
    }
}
