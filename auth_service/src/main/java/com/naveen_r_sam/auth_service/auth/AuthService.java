//	Copyright 2025 Naveen R
//
//		Licensed under the Apache License, Version 2.0 (the "License");
//		you may not use this file except in compliance with the License.
//		You may obtain a copy of the License at
//
//		http://www.apache.org/licenses/LICENSE-2.0
//
//		Unless required by applicable law or agreed to in writing, software
//		distributed under the License is distributed on an "AS IS" BASIS,
//		WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//		See the License for the specific language governing permissions and
//		limitations under the License.

package com.naveen_r_sam.auth_service.auth;

import com.naveen_r_sam.auth_service.common.MessageSender;
import com.naveen_r_sam.auth_service.dto.*;
import com.naveen_r_sam.auth_service.security.jwt.JWTService;
import com.naveen_r_sam.auth_service.model.Users;
import com.naveen_r_sam.auth_service.repo.UsersRepository;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.stereotype.Service;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;

import java.util.HashMap;
import java.util.Map;

@Service
public class AuthService implements IAuthService {
    private final UsersRepository _usersRepository;
    private final AuthenticationManager _authenticationManager;
    private final BCryptPasswordEncoder _passwordEncoder;
    private final JWTService _jwtService;
    private final MessageSender _messageSender;

    public AuthService(
            UsersRepository usersRepository,
            AuthenticationManager authenticationManager,
            JWTService jwtService,
            MessageSender messageSender
    ) {
        this._usersRepository = usersRepository;
        this._authenticationManager = authenticationManager;
        this._passwordEncoder = new BCryptPasswordEncoder(12);
        this._jwtService = jwtService;
        this._messageSender = messageSender;
    }

    public ResponseEntity<?> registerUser(SignUpDTO user) {
        if (user == null ||
                user.getUsername() == null ||
                user.getPassword() == null ||
                user.getFirstName() == null ||
                user.getLastName() == null ||
                user.getEmail() == null) {

            assert user != null;
            ErrorDataDTO err = getErrorDTOWithUserData(
                    user,
                    "user.signup",
                    400,
                    "Bad Request",
                    "POST"

            );
            this._messageSender.sendMessage(err.toString());
            return ResponseEntity.badRequest().body("All fields must be provided and non-null");
        }


        Users newUser = new Users();
        newUser.setFirstName(user.getFirstName());
        newUser.setLastName(user.getLastName());
        newUser.setEmail(user.getEmail());
        newUser.setUsername(user.getUsername());
        newUser.setPassword(_passwordEncoder.encode(user.getPassword()));

        Users savedUser = _usersRepository.save(newUser);
        if (savedUser.getId() == null) {
            ErrorDataDTO err = getErrorDTOWithUserData(
                    user,
                    "user.signup",
                    500,
                    "Internal Server Error",
                    "POST"
            );

            this._messageSender.sendMessage(err.toString());
            return ResponseEntity.internalServerError().body("User cannot be saved");
        }

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

    public ResponseEntity<?> authenticateUser(LoginRequestDTO user) {
        try {
            Authentication authentication = _authenticationManager.authenticate(
                    new UsernamePasswordAuthenticationToken(user.getUsername(), user.getPassword())
            );

            SecurityContextHolder.getContext().setAuthentication(authentication);
            if (!authentication.isAuthenticated()) {
                Map<String, Object> data = new HashMap<>();
                data.put("data", "Credentials not found (Wrong password)");

                ErrorDataDTO err = new ErrorDataDTO(
                        "user.login",
                        new ErrorDataDTO.ResponseData(
                                401,
                                "Unauthorized",
                                "POST",
                                data
                        )
                );
                this._messageSender.sendMessage(err.toString());
                return ResponseEntity.status(HttpStatus.UNAUTHORIZED).body("Credentials not found");
            }
            Users userData = _usersRepository.findByUsername(user.getUsername());
            String token = _jwtService.generateToken(userData);
            LoginResponseDTO response = new LoginResponseDTO(
                    token
            );
            return ResponseEntity.status(HttpStatus.OK).body(response);
        } catch (Exception ex) {
            Map<String, Object> data = new HashMap<>();
            data.put("data", "Credentials not found");

            ErrorDataDTO err = new ErrorDataDTO(
                    "user.login",
                    new ErrorDataDTO.ResponseData(
                            401,
                            "Unauthorized",
                            "POST",
                            data
                    )
            );
            this._messageSender.sendMessage(err.toString());
            return ResponseEntity.status(HttpStatus.UNAUTHORIZED).body("Credentials not found");
        }
    }

    private static ErrorDataDTO getErrorDTOWithUserData(SignUpDTO user, String pattern, int code, String status, String type) {
        Map<String, Object> data = new HashMap<>();
        Map<String, Object> userData = new HashMap<>();
        userData.put("firstName", user.getFirstName());
        userData.put("lastName", user.getLastName());
        userData.put("email", user.getEmail());
        userData.put("username", user.getUsername());
        data.put("data", userData);

        return new ErrorDataDTO(pattern,
                new ErrorDataDTO.ResponseData(
                        code,
                        status,
                        type,
                        data
                )
        );
    }
}
