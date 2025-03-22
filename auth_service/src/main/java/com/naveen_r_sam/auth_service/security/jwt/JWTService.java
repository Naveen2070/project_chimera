package com.naveen_r_sam.auth_service.security.jwt;

import com.naveen_r_sam.auth_service.model.Users;
import io.jsonwebtoken.Jwts;
import io.jsonwebtoken.security.Keys;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

import javax.crypto.SecretKey;
import java.nio.charset.StandardCharsets;
import java.util.*;

@Service
public class JWTService {
    @Value("${jwt.secret}")
    public String secretKey;

    public String generateToken(Users userData) {
        Map<String, Object> claims = new HashMap<>();
        claims.put("userId", userData.getId());
        claims.put("role", userData.getRole());

        return Jwts.builder().claims()
                .add(claims)
                .subject(userData.getUsername())
                .issuedAt(new Date(System.currentTimeMillis()))
                .expiration(new Date(System.currentTimeMillis() + 60 * 60 * 1000))
                .and()
                .signWith(generateKey(secretKey))
                .compact();
    }

    private SecretKey generateKey(String secret) {
        byte[] keyBytes = secret.getBytes(StandardCharsets.UTF_8);
        return Keys.hmacShaKeyFor(keyBytes);
    }
}