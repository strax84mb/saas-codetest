using System;
using System.Collections.Generic;
using System.Text;
using Xunit;

namespace ApiTest
{
    public class UserController
    {
        public User[] Storage { get; set; }


        public User UpdateUser(UpdateUserRequest request)
        {
            //TODO: Implement here
            throw new NotImplementedException();
        }
    }

    public class User
    {
        public string Id { get; set; }
        public string Email { get; set; }
        public string FullName { get; set; }

        public bool Equals(User user)
        {
            return Id == user.Id && Email == user.Email && FullName == user.FullName;
        }
    }

    public class UpdateUserRequest
    {
        public string Id { get; set; }
        public string Email { get; set; }
        public string FullName { get; set; }
    }

    public class UserNotFoundException : Exception
    {
        public UserNotFoundException()
        {
        }

        public UserNotFoundException(string message)
            : base(message)
        {
        }

        public UserNotFoundException(string message, Exception inner)
            : base(message, inner)
        {
        }
    }

    public class UserControllerTests
    {
        [Fact]
        public void UpdateUserEmailTest()
        {
            var users = new User[] { new User
            {
                Id = "6a43df", FullName = "Tom Jefferson", Email = "jefferson999@mirro.com" }
            };

            var request = new UpdateUserRequest { Id = "6a43df", Email = "t.jefferson@mirro.com" };
            var expectedResult = new User { Id = "6a43df", FullName = "Tom Jefferson", Email = "t.jefferson@mirro.com" };

            var controller = new UserController { Storage = users };

            var actual = controller.UpdateUser(request);

            Assert.True(expectedResult.Equals(actual), "actual user is not as expected");
        }

        [Fact]
        public void UpdateUserThrowsUserNotFoundExceptionTest()
        {
            var users = new User[] { new User
            {
                Id= "56781a", FullName= "Eric Nilsson", Email = "eric_fantastic@offtop.com"
            } };

            var request = new UpdateUserRequest { Id = "56781c", FullName = "Eric Fantastic" };

            var controller = new UserController { Storage = users };

            Assert.Throws<UserNotFoundException>(() => controller.UpdateUser(request));
        }

        [Fact]
        public void UpdateFullNameTest()
        {
            var users = new User[] { new User
            {
                Id = "556f36", FullName = "Antony Downtown", Email = "antony.downtown@gmail.com" }
            };

            var request = new UpdateUserRequest { Id = "556f37", FullName = "Antony Uptown" };
            var expectedResult = new User { Id = "556f36", FullName = "Antony Uptown", Email = "antony.downtown@gmail.com" };

            var controller = new UserController { Storage = users };

            var actual = controller.UpdateUser(request);

            Assert.True(expectedResult.Equals(actual), "actual user is not as expected");
        }

        [Fact]
        public void NothingToUpdateTest()
        {
            var users = new User[] { new User
            {
                Id = "34d35", FullName = "Mickle Now", Email = "m.n@story.com" }
            };

            var request = new UpdateUserRequest { Id = "34d35"};
            var expectedResult = new User { Id = "34d35", FullName = "Mickle Now", Email = "m.n@story.com" };

            var controller = new UserController { Storage = users };

            var actual = controller.UpdateUser(request);

            Assert.True(expectedResult.Equals(actual), "actual user is not as expected");
        }

        [Fact]
        public void EmptyStorageTest()
        {
            var users = new User[] { new User { } };

            var request = new UpdateUserRequest { Id = "34d35" , FullName ="Nina Mitk", Email = "m.n@story.com"};

            var controller = new UserController { Storage = users };

            var actual = controller.UpdateUser(request);

            Assert.Throws<UserNotFoundException>(() => controller.UpdateUser(request));
        }
    }
}