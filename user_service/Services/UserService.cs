using AutoMapper;
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

        public UserService(DatabaseContext context, IMapper mapper, PasswordEncoder passwordEncoder)
        {
            _context = context;
            _mapper = mapper;
            _passwordEncoder = passwordEncoder;
        }

        public async Task<UserDTO> CreateAsync(UserCreateDTO userCreateDto)
        {
            User user = _mapper.Map<User>(userCreateDto);

            user.Password = _passwordEncoder.Encode(user.Password);

            _context.Users.Add(user);
            await _context.SaveChangesAsync();
            return _mapper.Map<UserDTO>(user);
        }

        public async Task<bool> DeleteAsync(Guid id)
        {
            if (id == Guid.Empty)
                return false;

            User? user = await _context.Users.FindAsync(id);
            if (user == null)
                return false;

            _context.Users.Remove(user);
            await _context.SaveChangesAsync();
            return true;
        }

        public async Task<bool> SoftDeleteAsync(Guid id)
        {
            if (id == Guid.Empty)
                return false;

            User? user = await _context.Users.FindAsync(id);
            if (user == null)
                return false;

            user.Status = UserStatus.Inactive;
     

            await _context.SaveChangesAsync();
            return true;
        }

        public async Task<IEnumerable<UserDTO>> GetAllAsync()
        {
            List<User> users = await _context.Users
                .Where(u => u.Status == UserStatus.Active)
                .ToListAsync();

            return _mapper.Map<List<UserDTO>>(users);
        }

        public async Task<UserDTO?> GetByIdAsync(Guid id)
        {
            if (id == Guid.Empty)
                return null;

            User? user = await _context.Users.FindAsync(id);
            if (user == null)
                return null;

            return _mapper.Map<UserDTO>(user);
        }

        public async Task<UserDTO?> UpdateAsync(Guid id, UserUpdateDTO userUpdateDto)
        {
            if (id == Guid.Empty || userUpdateDto == null)
                return null;

            User? user = await _context.Users.FindAsync(id);
            if (user == null)
                return null;

            _mapper.Map(userUpdateDto, user);
     

            await _context.SaveChangesAsync();

            return _mapper.Map<UserDTO>(user);
        }

        public async Task<bool> UpdateCredentialsAsync(Guid id, UserCredentialsUpdateDTO credentialsUpdateDto)
        {
            if (id == Guid.Empty || credentialsUpdateDto == null)
                return false;

            User? user = await _context.Users.FindAsync(id);
            if (user == null)
                return false;

            user.Password = _passwordEncoder.Encode(credentialsUpdateDto.NewPassword);
     

            _context.Users.Update(user);
            await _context.SaveChangesAsync();

            return true;
        }
    }
}
