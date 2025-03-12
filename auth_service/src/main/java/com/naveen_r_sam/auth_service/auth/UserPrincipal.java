package com.naveen_r_sam.auth_service.auth;

import com.naveen_r_sam.auth_service.model.Users;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.core.userdetails.UserDetails;

import java.util.Collection;
import java.util.Collections;

public class UserPrincipal implements UserDetails {
    private final Users user;

    public UserPrincipal(Users user) {
        this.user = user;
    }

    @Override
    public Collection<? extends GrantedAuthority> getAuthorities() {
        return Collections.singleton(new SimpleGrantedAuthority("USER"));
    }

    @Override
    public String getPassword() {
        return user.getPassword();
    }

    @Override
    public String getUsername() {
        return user.getUsername();
    }

    @Override
    public boolean isAccountNonExpired() {
        return !user.getStatus().equals(0);
    }

    @Override
    public boolean isAccountNonLocked() {
        return !user.getStatus().equals(2);
    }

    @Override
    public boolean isCredentialsNonExpired() {
        return !user.getStatus().equals(3);
    }

    @Override
    public boolean isEnabled() {
        return user.getStatus().equals(1);
    }
}
