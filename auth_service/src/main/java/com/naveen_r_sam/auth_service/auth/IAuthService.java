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

import com.naveen_r_sam.auth_service.dto.LoginRequestDTO;
import com.naveen_r_sam.auth_service.dto.SignUpDTO;
import com.naveen_r_sam.auth_service.model.Users;
import org.springframework.http.ResponseEntity;

public interface IAuthService {
    /**
     * Register a user, given a {@link SignUpDTO} object.
     * @param user the user to register
     * @return a response entity containing a JSON object with
     * a single key, "message", containing a message describing the
     * result of the registration
     */
    ResponseEntity<?> registerUser(SignUpDTO user);
    /**
     * Authenticate a user, given a {@link LoginRequestDTO} object.
     * @param user the user to authenticate
     * @return a response entity containing a JSON object with
     * a single key, "message", describing the result of the authentication,
     * and possibly additional authentication information if successful
     */
    ResponseEntity<?> authenticateUser(LoginRequestDTO user);
}
