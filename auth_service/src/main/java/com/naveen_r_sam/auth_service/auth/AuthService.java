package com.naveen_r_sam.auth_service.auth;

import com.naveen_r_sam.auth_service.dto.LoginResponseDTO;
import com.naveen_r_sam.auth_service.security.jwt.JWTService;
import com.naveen_r_sam.auth_service.model.Users;
import com.naveen_r_sam.auth_service.dto.UsersDTO;
import com.naveen_r_sam.auth_service.repo.UsersRepository;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.stereotype.Service;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;

@Service
public class AuthService implements IAuthService {
    private final UsersRepository _usersRepository;
    private final AuthenticationManager _authenticationManager;
    private final BCryptPasswordEncoder _passwordEncoder;
    private final JWTService _jwtService;

    public AuthService(UsersRepository usersRepository, AuthenticationManager authenticationManager, JWTService jwtService) {
        this._usersRepository = usersRepository;
        this._authenticationManager = authenticationManager;
        this._passwordEncoder = new BCryptPasswordEncoder(12);
        this._jwtService = jwtService;
    }

    public ResponseEntity<?> registerUser(Users user) {
        if (user == null) {
            return ResponseEntity.badRequest().body("User cannot be null");
        }

        user.setPassword(_passwordEncoder.encode(user.getPassword()));

        Users savedUser = _usersRepository.save(user);
        UsersDTO response = new UsersDTO(
                savedUser.getId(),
                savedUser.getFirstName(),
                savedUser.getLastName(),
                savedUser.getEmail(),
                savedUser.getUsername(),
                savedUser.getRole(),
                savedUser.getStatus(),
                savedUser.getCreatedOn(),
                savedUser.getUpdatedOn()
        );

        return ResponseEntity.status(HttpStatus.CREATED).body(response);
    }

    public ResponseEntity<?> authenticateUser(Users user) {

        try {
            Authentication authentication = _authenticationManager.authenticate(
                    new UsernamePasswordAuthenticationToken(user.getUsername(), user.getPassword())
            );

            SecurityContextHolder.getContext().setAuthentication(authentication);
            if (!authentication.isAuthenticated()) {
                return ResponseEntity.status(HttpStatus.UNAUTHORIZED).body("Credentials not found");
            }
            String token = _jwtService.generateToken(user.getUsername());
            LoginResponseDTO response = new LoginResponseDTO(
                    token
            );
            return ResponseEntity.status(HttpStatus.OK).body(response);
        } catch (Exception ex) {
            return ResponseEntity.status(HttpStatus.UNAUTHORIZED).body("Credentials not found");
        }
    }
}
