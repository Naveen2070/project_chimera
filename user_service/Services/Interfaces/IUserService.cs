using user_service.Model;

namespace user_service.Services.Interfaces
{
    public interface IUserService
    {
        Task<IEnumerable<UserDTO>> GetAllAsync();
        Task<UserDTO?> GetByIdAsync(Guid id);
        Task<UserDTO> CreateAsync(UserCreateDTO userCreateDto);
        Task<UserDTO?> UpdateAsync(Guid id, UserUpdateDTO userUpdateDto);
        Task<bool> DeleteAsync(Guid id);
        Task<bool> SoftDeleteAsync(Guid id);
        Task<bool> UpdateCredentialsAsync(Guid id, UserCredentialsUpdateDTO credentialsUpdateDto);
    }
}
