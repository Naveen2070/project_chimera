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
using AutoMapper;
using Microsoft.AspNetCore.Http.HttpResults;
using Microsoft.EntityFrameworkCore;
using user_service.Database;
using user_service.Entity.AuthService.Model;
using user_service.Model;
using user_service.Services.Interfaces;
using user_service.Utils;
using user_service.Utils.Enums;

namespace user_service.Services
{
    public class UserService : IUserService
    {
        private readonly DatabaseContext _context;
        private readonly IMapper _mapper;
        private readonly PasswordEncoder _passwordEncoder;
        private readonly IRMQPublisherService _publisher;

        public UserService(DatabaseContext context, IMapper mapper, PasswordEncoder passwordEncoder, IRMQPublisherService publisher)
        {
            _context = context;
            _mapper = mapper;
            _passwordEncoder = passwordEncoder;
            _publisher = publisher;
        }

        public async Task<UserDTO> CreateAsync(UserCreateDTO userCreateDto)
        {
            User user;
            try
            {
                user = _mapper.Map<User>(userCreateDto);

                user.Password = _passwordEncoder.Encode(user.Password);

                _context.Users.Add(user);

                await _context.SaveChangesAsync();
            }
            catch (Exception e)
            {
                if (e.InnerException != null && e.InnerException.Message.Contains("unique constraint", StringComparison.OrdinalIgnoreCase))
                {
                    var error = ErrorDTO.CreateErrorDTO(
                        pattern: "user.create",
                         code: 400,
                         type: "Post",
                         status: "Bad Request",
                         data: new
                         {
                             message = "User with email already exists."
                         }
                     );
                    await _publisher.SendMessageAsync(error.ToString());
                    throw new Exception("User with email already exists.");
                }
                var generalError = ErrorDTO.CreateErrorDTO(
                        pattern: "user.create",
                         code: 400,
                         type: "Post",
                         status: "Bad Request",
                         data: new
                         {
                             message = e.ToString()
                         }
                     );
                await _publisher.SendMessageAsync(generalError.ToString());
                throw new Exception(e.Message);
            }
            return _mapper.Map<UserDTO>(user);
        }

        public async Task<bool> DeleteAsync(Guid id)
        {
            if (id == Guid.Empty)
                return false;

            try
            {
                var user = await _context.Users.FindAsync(id);

                if (user == null)
                {
                    var notFound = ErrorDTO.CreateErrorDTO(
                        pattern: "user.delete",
                        code: 404,
                        type: "Delete",
                        status: "Not Found",
                        data: new
                        {
                            message = "User with ID " + id + " not found."
                        }
                    );

                    await _publisher.SendMessageAsync(notFound.ToString());
                    return false;
                }

                _context.Users.Remove(user);
                await _context.SaveChangesAsync();
                return true;
            }
            catch (Exception e)
            {
                var error = ErrorDTO.CreateErrorDTO(
                    pattern: "user.delete",
                    code: 500,
                    type: "Delete",
                    status: "Internal Server Error",
                    data: new
                    {
                        message = e.ToString()
                    }
                );

                await _publisher.SendMessageAsync(error.ToString());
                throw; 
            }
        }


        public async Task<bool> SoftDeleteAsync(Guid id)
        {
            if (id == Guid.Empty)
                return false;

            try
            {
                var user = await _context.Users.FindAsync(id);
                if (user == null)
                {
                    var notFound = ErrorDTO.CreateErrorDTO(
                        pattern: "user.softdelete",
                        code: 404,
                        type: "SoftDelete",
                        status: "Not Found",
                        data: new { message = "User with ID " + id + " not found." }
                    );
                    await _publisher.SendMessageAsync(notFound.ToString());
                    return false;
                }

                user.Status = UserStatus.Inactive;
                await _context.SaveChangesAsync();

                var success = ErrorDTO.CreateErrorDTO(
                    pattern: "user.softdelete",
                    code: 200,
                    type: "SoftDelete",
                    status: "Success",
                    data: new
                    {
                        message = $"User with ID {id} was soft-deleted.",
                        userId = id
                    }
                );
                await _publisher.SendMessageAsync(success.ToString());

                return true;
            }
            catch (Exception ex)
            {
                var error = ErrorDTO.CreateErrorDTO(
                    pattern: "user.softdelete",
                    code: 500,
                    type: "SoftDelete",
                    status: "Internal Server Error",
                    data: new { message = ex.ToString() }
                );
                await _publisher.SendMessageAsync(error.ToString());
                throw;
            }
        }


        public async Task<IEnumerable<UserDTO>> GetAllAsync()
        {
            try
            {
                var users = await _context.Users
                    .Where(u => u.Status == UserStatus.Active)
                    .ToListAsync();

                var userDTOs = _mapper.Map<List<UserDTO>>(users);

                var successMessage = ErrorDTO.CreateErrorDTO(
                    pattern: "user.getall",
                    code: 200,
                    type: "GetAll",
                    status: "Success",
                    data: new
                    {
                        message = "Fetched all active users.",
                        count = userDTOs.Count
                    }
                );
                await _publisher.SendMessageAsync(successMessage.ToString());

                return userDTOs;
            }
            catch (Exception ex)
            {
                var errorMessage = ErrorDTO.CreateErrorDTO(
                    pattern: "user.getall",
                    code: 500,
                    type: "GetAll",
                    status: "Internal Server Error",
                    data: new
                    {
                        message = ex.ToString()
                    }
                );
                await _publisher.SendMessageAsync(errorMessage.ToString());
                throw;
            }
        }


        public async Task<UserDTO?> GetByIdAsync(Guid id)
        {
            if (id == Guid.Empty)
                return null;

            try
            {
                var user = await _context.Users.FindAsync(id);
                if (user == null || user.Status != UserStatus.Active)
                {
                    var notFound = ErrorDTO.CreateErrorDTO(
                        pattern: "user.getbyid",
                        code: 404,
                        type: "GetById",
                        status: "Not Found",
                        data: new { message = "User " + id + " not found." }
                    );
                    await _publisher.SendMessageAsync(notFound.ToString());
                    return null;
                }

                var success = ErrorDTO.CreateErrorDTO(
                    pattern: "user.getbyid",
                    code: 200,
                    type: "GetById",
                    status: "Success",
                    data: new { message = "User retrieved.", userId = id }
                );
                await _publisher.SendMessageAsync(success.ToString());

                return _mapper.Map<UserDTO>(user);
            }
            catch (Exception ex)
            {
                var error = ErrorDTO.CreateErrorDTO(
                    pattern: "user.getbyid",
                    code: 500,
                    type: "GetById",
                    status: "Internal Server Error",
                    data: new { message = ex.ToString() }
                );
                await _publisher.SendMessageAsync(error.ToString());
                throw;
            }
        }

        public async Task<UserDTO?> UpdateAsync(Guid id, UserUpdateDTO userUpdateDto)
        {
            if (id == Guid.Empty || userUpdateDto == null)
                return null;

            try
            {
                var user = await _context.Users.FindAsync(id);
                if (user == null || user.Status != UserStatus.Active)
                {
                    var notFound = ErrorDTO.CreateErrorDTO(
                        pattern: "user.update",
                        code: 404,
                        type: "Update",
                        status: "Not Found",
                        data: new { message = "User with ID " + id + " not found." }
                    );
                    await _publisher.SendMessageAsync(notFound.ToString());
                    return null;
                }

                _mapper.Map(userUpdateDto, user);
                await _context.SaveChangesAsync();

                var success = ErrorDTO.CreateErrorDTO(
                    pattern: "user.update",
                    code: 200,
                    type: "Update",
                    status: "Success",
                    data: new { message = "User updated.", userId = id }
                );
                await _publisher.SendMessageAsync(success.ToString());

                return _mapper.Map<UserDTO>(user);
            }
            catch (Exception ex)
            {
                var error = ErrorDTO.CreateErrorDTO(
                    pattern: "user.update",
                    code: 500,
                    type: "Update",
                    status: "Internal Server Error",
                    data: new { message = ex.ToString() }
                );
                await _publisher.SendMessageAsync(error.ToString());
                throw;
            }
        }


        public async Task<bool> UpdateCredentialsAsync(Guid id, UserCredentialsUpdateDTO credentialsUpdateDto)
        {
            if (id == Guid.Empty || credentialsUpdateDto == null)
                return false;

            try
            {
                var user = await _context.Users.FindAsync(id);
                if (user == null || user.Status != UserStatus.Active)
                {
                    var notFound = ErrorDTO.CreateErrorDTO(
                        pattern: "user.updatecredentials",
                        code: 404,
                        type: "UpdateCredentials",
                        status: "Not Found",
                        data: new { message = "User not with ID " + id + " found." }
                    );
                    await _publisher.SendMessageAsync(notFound.ToString());
                    return false;
                }

                user.Password = _passwordEncoder.Encode(credentialsUpdateDto.NewPassword);
                await _context.SaveChangesAsync();

                var success = ErrorDTO.CreateErrorDTO(
                    pattern: "user.updatecredentials",
                    code: 200,
                    type: "UpdateCredentials",
                    status: "Success",
                    data: new { message = "User credentials updated.", userId = id }
                );
                await _publisher.SendMessageAsync(success.ToString());

                return true;
            }
            catch (Exception ex)
            {
                var error = ErrorDTO.CreateErrorDTO(
                    pattern: "user.updatecredentials",
                    code: 500,
                    type: "UpdateCredentials",
                    status: "Internal Server Error",
                    data: new { message = ex.ToString() }
                );
                await _publisher.SendMessageAsync(error.ToString());
                throw;
            }
        }

    }
}
