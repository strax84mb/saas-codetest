using System;
using Xunit;

namespace Tests
{
    public class DeepestLetterTests
    {
		public static char DeepestLetter(string str)
		{
			throw new NotImplementedException();
		}

		[Fact]
		public void DeepestLetterSimpleTest()
        {
			var input = "a(b)c";
			var expected = 'b';
			var actual = DeepestLetter(input);

			Assert.Equal(expected, actual);
        }

		[Fact]
		public void DeepestLetterComplexTest()
		{
			var input = "((a))(((M)))(c)(D)(e)(((f))(((G))))h(i)";
			var expected = 'g';
			var actual = DeepestLetter(input);

			Assert.Equal(expected, actual);
		}

		[Fact]
		public void DeepestLetterFirstMetTest()
		{
			var input = "((A)(b)c";
			var expected = 'c';
			var actual = DeepestLetter(input);

			Assert.Equal(expected, actual);
		}

		[Fact]
		public void DeepestLetterOneNestedTest()
		{
			var input = "(a)((G)c)";
			var expected = 'g';
			var actual = DeepestLetter(input);

			Assert.Equal(expected, actual);
		}

		[Fact]
		public void DeepestLetterNumberTest()
		{
			var input = "(8)";
			var expected = '?';
			var actual = DeepestLetter(input);

			Assert.Equal(expected, actual);
		}

		[Fact]
		public void DeepestLetterSymbolTest()
		{
			var input = "(!)";
			var expected = 'a';
			var actual = DeepestLetter(input);

			Assert.Equal(expected, actual);
		}
	}
}
