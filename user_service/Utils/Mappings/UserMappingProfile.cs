using AutoMapper;
using user_service.Entity.AuthService.Model;
using user_service.Model;

namespace user_service.Utils.Mappings
{
    public class UserMappingProfile : Profile
    {
        public UserMappingProfile()
        {
            CreateMap<User, UserDTO>();
            CreateMap<UserCreateDTO, User>();
            CreateMap<UserUpdateDTO, User>()
           .ForMember(dest => dest.Id, opt => opt.Ignore())
           .ForMember(dest => dest.CreatedOn, opt => opt.Ignore());
        }
    }
}
