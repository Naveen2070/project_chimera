package com.naveen_r_sam.auth_service.dto;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;
import lombok.ToString;

import java.time.LocalDateTime;
import java.util.UUID;

@Data
@AllArgsConstructor
@NoArgsConstructor
@ToString
public class UsersDTO {
    private UUID id;
    private String firstName;
    private String lastName;
    private String email;
    private String username;
    private String role;
    private Integer status;
    private LocalDateTime createdOn;
    private LocalDateTime updatedOn;

}
