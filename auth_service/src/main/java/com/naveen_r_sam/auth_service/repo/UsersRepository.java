package com.naveen_r_sam.auth_service.repo;

import com.naveen_r_sam.auth_service.model.Users;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.UUID;

public interface UsersRepository extends JpaRepository<Users, UUID> {
    Users findByUsername(String username);
}
